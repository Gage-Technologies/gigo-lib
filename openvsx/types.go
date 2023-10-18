package openvsx

type Extension struct {
	NamespaceURL       string             `json:"namespaceUrl"`
	ReviewsURL         string             `json:"reviewsUrl"`
	Files              map[string]string  `json:"files"`
	Name               string             `json:"name"`
	Namespace          string             `json:"namespace"`
	TargetPlatform     string             `json:"targetPlatform"`
	Version            string             `json:"version"`
	PreRelease         bool               `json:"preRelease"`
	PublishedBy        Publisher          `json:"publishedBy"`
	Verified           bool               `json:"verified"`
	UnrelatedPublisher bool               `json:"unrelatedPublisher"`
	NamespaceAccess    string             `json:"namespaceAccess"`
	AllVersions        map[string]string  `json:"allVersions"`
	AllVersionsURL     string             `json:"allVersionsUrl"`
	AverageRating      float64            `json:"averageRating"`
	DownloadCount      int                `json:"downloadCount"`
	ReviewCount        int                `json:"reviewCount"`
	VersionAlias       []string           `json:"versionAlias"`
	Timestamp          string             `json:"timestamp"`
	Preview            bool               `json:"preview"`
	DisplayName        string             `json:"displayName"`
	Description        string             `json:"description"`
	Engines            map[string]string  `json:"engines"`
	Categories         []string           `json:"categories"`
	ExtensionKind      []string           `json:"extensionKind"`
	Tags               []string           `json:"tags"`
	License            string             `json:"license"`
	Homepage           string             `json:"homepage"`
	Repository         string             `json:"repository"`
	SponsorLink        string             `json:"sponsorLink"`
	Bugs               string             `json:"bugs"`
	GalleryColor       string             `json:"galleryColor"`
	GalleryTheme       string             `json:"galleryTheme"`
	LocalizedLanguages []string           `json:"localizedLanguages"`
	QnA                string             `json:"qna"`
	Dependencies       []interface{}      `json:"dependencies"`
	BundledExtensions  []BundledExtension `json:"bundledExtensions"`
	Downloads          map[string]string  `json:"downloads"`
}

type Publisher struct {
	LoginName string `json:"loginName"`
	FullName  string `json:"fullName"`
	AvatarURL string `json:"avatarUrl"`
	Homepage  string `json:"homepage"`
	Provider  string `json:"provider"`
}

type BundledExtension struct {
	URL       string `json:"url"`
	Namespace string `json:"namespace"`
	Extension string `json:"extension"`
}
