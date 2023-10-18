package spice

import (
	"context"
	"database/sql"
	"fmt"
	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	ti "github.com/gage-technologies/gigo-lib/db"
	"google.golang.org/grpc"
	"io"
	"time"
)

type Database struct {
	DB    *authzed.Client
	Token string
	tiDB  *ti.Database
}

// CreateDatabase
// Create a new database instance for SpiceDB
// Args
//
//	host 		- string, host address for the SpiceDB instance
//	port 		- string, port that the SpiceDB instance is available at
//	token 		- string, bearer token for authentication with the SpiceDB instance
//	tls 		- bool, whether to use TLS encryption to secure connection with SpiceDB
//	cert 		- string, path to certificate file for use in TLS connection
//
// Returns
//
//	out 		- *Database, new instance of the SpiceDB database connection
func CreateDatabase(host string, port string, token string, tls bool, cert string, tiDB *ti.Database) (*Database, error) {
	// create slice to hold dial options to permit variable dial options
	var dialOptions []grpc.DialOption

	// assemble dial options depending on whether the
	if tls {
		dialOptions = []grpc.DialOption{
			grpcutil.WithCustomCerts(cert, true),
			grpcutil.WithBearerToken(token),
		}
	} else {
		dialOptions = []grpc.DialOption{
			grpc.WithInsecure(),
			grpcutil.WithInsecureBearerToken(token),
		}
	}

	// create a new client connected to the spicedb instance
	client, err := authzed.NewClient(
		fmt.Sprintf("%s:%s", host, port),
		dialOptions...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to spicedb: %v", err)
	}

	//// conditionally attempt to initialize the zookie tracker table
	//if tiDB != nil {
	//	_, err := tiDB.DB.Exec(
	//		"create table if not exists zookies(" +
	//			"r_id varchar(36) not null, " +
	//			"s_id varchar(36) not null, " +
	//			"z varchar(64) not null, " +
	//			"time datetime not null, " +
	//			"primary key (r_id, s_id)" +
	//			")",
	//	)
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to initliaze the zookie table: %v", err)
	//	}
	//}

	return &Database{
		DB:    client,
		Token: token,
		tiDB:  tiDB,
	}, nil
}

// updateZookieTable
// Internal function to automatically update the zookie table when performing relation ops
// Args
//
//	relations 			- []Relation, relations that were updated in SpiceDB
//	zookie 				- string, the Zookie (ZedToken) returned by the relation op
func (db *Database) updateZookieTable(relations []Relation, zookie string) error {
	// return silently if there is no tidb instance configured
	if db.tiDB == nil {
		return nil
	}

	// create base insert statement
	insertStatement := "insert into zookies(r_id, s_id, z, time) values %s on duplicate key update z = values(z), time = values(time)"
	// create string to hold variable sets for the insertion ops
	variableSets := ""
	// create slice to hold parameters of the variable sets
	params := make([]interface{}, 0)

	// iterate over relations updating the variables sets and parameters
	for i, relation := range relations {
		// conditionally add a comma and space to the variable sets
		if i != 0 {
			variableSets += ", "
		}
		// add a variable set to the string
		variableSets += "(?, ?, ?, ?)"
		// append the values for the variable set to the params
		params = append(params, relation.ResourceID, relation.SubjectID, zookie, time.Now())
	}

	// format the full statement and execute the update operation
	_, err := db.tiDB.DB.Exec(fmt.Sprintf(insertStatement, variableSets), params...)
	if err != nil {
		return fmt.Errorf("failed to update zookie tracker table: %v\n    query: %s\n    params: %v", err, fmt.Sprintf(insertStatement, variableSets), params)
	}

	return nil
}

// getZookie
// Attempts to retrieve an existing zookie from the zookie tracker table
// Args
//
//	resourceId 			- string, id of the resource associated with the zookie
//	subjectId 			- string, id iof the subject associated with teh zookie
//
// Returns
//
//	zookie 				- string, zookie found in the zookie tracker table
func (db *Database) getZookie(resourceId, subjectId string) (string, error) {
	// silently return an empty string if no tidb instance was configured
	if db.tiDB == nil {
		return "", nil
	}

	// ensure that only one of the parameters was passed
	if resourceId == "" && subjectId == "" {
		return "", fmt.Errorf("at leat one parameter should be passed for getZookie")
	}

	// create default command to retrieve zookie
	query := "select z from zookies where %s order by time desc limit 1"
	// create filter string  and params for query
	filter := ""
	params := make([]interface{}, 0)

	// conditionally use resource
	if resourceId != "" {
		filter += "r_id = ?"
		params = append(params, resourceId)
	}

	// conditionally use subject
	if subjectId != "" {
		// conditionally add `or` to filter
		if len(params) > 0 {
			filter += " or "
		}
		filter += "s_id = ?"
		params = append(params, subjectId)
	}

	// attempt to retrieve zookie from database
	row := db.tiDB.DB.QueryRow(fmt.Sprintf(query, filter), params...)

	// create an empty string to load zookie into
	var zookie string

	// attempt to load zookie from row
	err := row.Scan(&zookie)
	if err != nil {
		// return an empty string if zookie doesn't exist
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to load zookie for resource %q and subject %q: %v", resourceId, subjectId, err)
	}

	return zookie, nil
}

// CreateRelations
// Wrapper function to simplify interface for the creation of relations in SpiceDB
// Args
//
//	relations 			- ...Relation, variable list of relations that will be created
//
// Returns
//
//	zookie 				- string, Zookie (ZedToken) to track cache consistency
func (db *Database) CreateRelations(relations ...Relation) (string, error) {
	// create slice to hold relations in the native GRPC format
	nativeRelations := make([]*pb.RelationshipUpdate, len(relations))

	// iterate over passed permissions reformatting them into GRPC format
	for i, r := range relations {
		nativeRelations[i] = &pb.RelationshipUpdate{
			Operation: pb.RelationshipUpdate_OPERATION_TOUCH,
			Relationship: &pb.Relationship{
				Resource: &pb.ObjectReference{
					ObjectType: string(r.ResourceType),
					ObjectId:   r.ResourceID,
				},
				Relation: string(r.Relation),
				Subject: &pb.SubjectReference{
					Object: &pb.ObjectReference{
						ObjectType: string(r.SubjectType),
						ObjectId:   r.SubjectID,
					},
				},
			},
		}
	}

	// execute relation insertions
	res, err := db.DB.WriteRelationships(context.TODO(), &pb.WriteRelationshipsRequest{Updates: nativeRelations})
	if err != nil {
		return "", fmt.Errorf("failed to update relations in spicedb: %v", err)
	}

	// perform zookie table update
	err = db.updateZookieTable(relations, res.GetWrittenAt().GetToken())
	if err != nil {
		return "", fmt.Errorf("failed to update zookie table: %v", err)
	}

	return res.GetWrittenAt().GetToken(), nil
}

// CheckPermission
// Wrapper function to simplify interface for the checking of permissions in SpiceDB
// Args
//
//	permission 			- Permission, permission that will be checked for in SpiceDB
//	zookie 				- string, optional Zookie (ZedToken) to track cache consistency (pass full for full consistency)
//
// Returns
//
//	permissionCheck 	- CheckPermission, the result of the permission check represented in the native type
//	zookie 				- string, Zookie (ZedToken) returned by the permission check for tracking consistency
func (db *Database) CheckPermission(permission Permission, zookie string) (CheckPermissionType, string, error) {
	// attempt to retrieve zookie if no zookie was passed
	if zookie == "" {
		// execute zookie retrieval
		z, err := db.getZookie(permission.ResourceID, permission.SubjectID)
		if err != nil {
			return "", "", fmt.Errorf("failed to retrieve zookie for permission check: %v", err)
		}

		// update zookie
		zookie = z
	}

	// create default consistency check of full consistency
	consistencyCheck := &pb.Consistency{
		Requirement: &pb.Consistency_FullyConsistent{
			FullyConsistent: true,
		},
	}

	// conditional consistency check
	if zookie != "" && zookie != "full" {
		// create consistency check
		consistencyCheck = &pb.Consistency{
			Requirement: &pb.Consistency_AtLeastAsFresh{
				AtLeastAsFresh: &pb.ZedToken{
					Token: zookie,
				},
			},
		}
	}

	// execute permission check using GRPC client
	res, err := db.DB.CheckPermission(context.TODO(), &pb.CheckPermissionRequest{
		Consistency: consistencyCheck,
		Resource: &pb.ObjectReference{
			ObjectType: string(permission.ResourceType),
			ObjectId:   permission.ResourceID,
		},
		Permission: string(permission.Permission),
		Subject: &pb.SubjectReference{
			Object: &pb.ObjectReference{
				ObjectType: string(permission.SubjectType),
				ObjectId:   permission.SubjectID,
			},
		},
	})
	if err != nil {
		return "", "", err
	}

	// perform zookie table update
	err = db.updateZookieTable([]Relation{{
		ResourceID: permission.ResourceID,
		SubjectID:  permission.SubjectID,
	}}, res.GetCheckedAt().GetToken())
	if err != nil {
		return "", "", fmt.Errorf("failed to update zookie table: %v", err)
	}

	// convert check permission result to native status
	switch res.Permissionship {
	case pb.CheckPermissionResponse_PERMISSIONSHIP_UNSPECIFIED:
		return UnspecifiedPermission, res.GetCheckedAt().GetToken(), nil
	case pb.CheckPermissionResponse_PERMISSIONSHIP_NO_PERMISSION:
		return NoPermission, res.GetCheckedAt().GetToken(), nil
	case pb.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION:
		return HasPermission, res.GetCheckedAt().GetToken(), nil
	default:
		return "", res.GetCheckedAt().GetToken(), fmt.Errorf("unexpected permission result returned from spicedb %q", res.Permissionship.String())
	}
}

// DeleteRelation
// Wrapper function to simplify interface for the deleting of relations in SpiceDB
// Args
//
//	relation 			- Relation, relation that will be deleted from SpiceDB
//
// Returns
//
//	zookie 				- string, Zookie (ZedToken) returned by the relation deletion for tracking consistency
func (db *Database) DeleteRelation(relation Relation) (string, error) {
	// create base request with default options (optional fields here are strings; they can be filled even if they were not passed)
	req := pb.RelationshipFilter{
		ResourceType:       string(relation.ResourceType),
		OptionalResourceId: relation.ResourceID,
		OptionalRelation:   string(relation.Relation),
	}

	// conditionally include a subject filter if values were passed
	if relation.SubjectType != "" {
		// include the subject information
		req.OptionalSubjectFilter = &pb.SubjectFilter{
			SubjectType:       string(relation.SubjectType),
			OptionalSubjectId: relation.SubjectID,
		}
	}

	// execute permission removal call
	res, err := db.DB.DeleteRelationships(context.TODO(), &pb.DeleteRelationshipsRequest{RelationshipFilter: &req})
	if err != nil {
		return "", nil
	}

	// perform zookie table update
	err = db.updateZookieTable([]Relation{relation}, res.GetDeletedAt().GetToken())
	if err != nil {
		return "", fmt.Errorf("failed to update zookie table: %v", err)
	}

	return res.GetDeletedAt().GetToken(), err
}

// GetPermissions
// Wrapper function to simplify interface for the retrieval of resources a user has access to
// Args
//
//	resourceType 		- ObjectType, type of the object that we will be bound the permissions to
//	subjectType 		- ObjectType, type of the subject that we are checking permission for
//	subjectId 			- string, id of the subject that we are checking for
//	permission 			- PermissionType, type of permission that we are checking for
//	zookie 				- string, optional Zookie (ZedToken) to track cache consistency (pass full for full consistency)
//
// Returns
//
//	resources 			- []string, slice containing the ids of each resource that the subject has the specified permission for
func (db *Database) GetPermissions(resourceType ObjectType, subjectType ObjectType, subjectId string,
	permission PermissionType, zookie string) ([]string, error) {
	// attempt to retrieve zookie if no zookie was passed
	if zookie == "" {
		// execute zookie retrieval
		z, err := db.getZookie("", subjectId)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve zookie for permission check: %v", err)
		}

		// update zookie
		zookie = z
	}

	// create default consistency check of full consistency
	consistencyCheck := &pb.Consistency{
		Requirement: &pb.Consistency_FullyConsistent{
			FullyConsistent: true,
		},
	}

	// conditional consistency check
	if zookie != "" && zookie != "full" {
		// create consistency check
		consistencyCheck = &pb.Consistency{
			Requirement: &pb.Consistency_AtLeastAsFresh{
				AtLeastAsFresh: &pb.ZedToken{
					Token: zookie,
				},
			},
		}
	}

	// execute lookup operation
	stream, err := db.DB.LookupResources(context.TODO(), &pb.LookupResourcesRequest{
		Consistency:        consistencyCheck,
		ResourceObjectType: string(resourceType),
		Permission:         string(permission),
		Subject: &pb.SubjectReference{
			Object: &pb.ObjectReference{
				ObjectType: string(subjectType),
				ObjectId:   subjectId,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute lookup resource operation: %v", err)
	}

	// create slice to hold resource ids
	resources := make([]string, 0)

	// iterate over stream reading the resources that the subject has permissions to
	for {
		// retrieve next resource from stream
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		// handle error
		if err != nil {
			return nil, fmt.Errorf("failed to read from stream while loading lookup resources: %v", err)
		}

		// retrieve id from resource
		resources = append(resources, res.ResourceObjectId)
	}

	return resources, nil
}

// GetResourceRelations
// Wrapper function to simplify interface for the retrieval of relations for a given resource
// Args
//
//	resourceType 		- ObjectType, type of the object that we will be bound the relations to
//	resourceId 			- string, id of the resource that we are checking for
//	relation 			- RelationType, type of relation that we are checking for
//	zookie 				- string, optional Zookie (ZedToken) to track cache consistency (pass full for full consistency)
//
// Returns
//
//	relations 			- []Relation, slice containing the relations of the resource was specified
func (db *Database) GetResourceRelations(resourceType ObjectType, resourceId string, relation RelationType,
	zookie string) ([]Relation, error) {
	// attempt to retrieve zookie if no zookie was passed
	if zookie == "" {
		// execute zookie retrieval
		z, err := db.getZookie("", resourceId)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve zookie for permission check: %v", err)
		}

		// update zookie
		zookie = z
	}

	// create default consistency check of full consistency
	consistencyCheck := &pb.Consistency{
		Requirement: &pb.Consistency_FullyConsistent{
			FullyConsistent: true,
		},
	}

	// conditional consistency check
	if zookie != "" && zookie != "full" {
		// create consistency check
		consistencyCheck = &pb.Consistency{
			Requirement: &pb.Consistency_AtLeastAsFresh{
				AtLeastAsFresh: &pb.ZedToken{
					Token: zookie,
				},
			},
		}
	}

	// execute relation retrieval operation
	stream, err := db.DB.ReadRelationships(context.TODO(), &pb.ReadRelationshipsRequest{
		Consistency: consistencyCheck,
		RelationshipFilter: &pb.RelationshipFilter{
			ResourceType:       string(resourceType),
			OptionalResourceId: resourceId,
			OptionalRelation:   string(relation),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute read relationships operation: %v", err)
	}

	// create slice to hold relations
	relations := make([]Relation, 0)

	// iterate over stream reading the relations that the resource is related to
	for {
		// retrieve next relation from stream
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		// handle error
		if err != nil {
			return nil, fmt.Errorf("failed to read from stream while loading read realtionshipd: %v", err)
		}

		// retrieve relation from response
		relations = append(relations, Relation{
			ResourceType: ObjectType(res.GetRelationship().GetResource().GetObjectType()),
			ResourceID:   res.GetRelationship().GetResource().GetObjectId(),
			Relation:     RelationType(res.GetRelationship().GetRelation()),
			SubjectType:  ObjectType(res.GetRelationship().GetSubject().GetObject().GetObjectType()),
			SubjectID:    res.GetRelationship().GetSubject().GetObject().GetObjectId(),
		})
	}

	return relations, nil
}
