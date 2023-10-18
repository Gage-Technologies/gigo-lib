import hashlib
import io
import os
import json
import argparse
import tempfile
import traceback as tb
from typing.io import IO
from urllib import request
from subprocess import Popen, PIPE

from typing import Tuple, Optional

# load arguments
parser = argparse.ArgumentParser()
parser.add_argument('--workspace-id', type=str, required=True)
args = parser.parse_args()

WORKING_DIRECTORY = "<working_directory>"

CONTAINER_COMPOSE = """
<containers>
"""

SHELL_EXECUTIONS = """[]"""

USE_VSCODE = True

VSCODE_EXTENSIONS = []


def execute_command(command: str, log_file: Optional[str] = None) -> Tuple[int, str, str]:
    """
    execute_command executes a command and returns the status, stdout and stderr
    :param command: arbitrary shell command to be executed
    :param log_file: optional log file to write the output to
    :return: status, stdout and stderr of the command
    """
    if log_file is None:
        h = hashlib.sha3_256(command.encode()).hexdigest()
        log_file = tempfile.mktemp(prefix=f"gigo-ws-init-cmd-{h[:7]}")

    with open(log_file, "w+b") as out_pipe:
        with open(log_file + ".err", "w+b") as err_pipe:
            p = Popen(command, shell=True, stdout=out_pipe, stderr=err_pipe)
            status = p.wait()

            out_pipe.flush()
            err_pipe.flush()

            # get file sizes so that we can read only the last 64Kib for the log return
            out_size = out_pipe.tell()
            err_size = err_pipe.tell()

            # read at most the last 64Kib of the logs
            out_pipe.seek(max(out_size-65536, 0), 0)
            err_pipe.seek(max(err_size-65536, 0), 0)

            out_buf = out_pipe.read()
            err_buf = err_pipe.read()

    return status, out_buf.decode('utf-8'), err_buf.decode('utf-8')


def download_extension(_id: int, secret: str):
    """
    download_extension downloads the gigo-developer extension from the gigo remote server
    :param _id: the id of the workspace
    :param secret: the secret to download the extension from
    """
    # format request data
    request_data = json.dumps({
        "workspace_id": _id,
        "secret": secret,
    })

    # format POST request
    req = request.Request(
        "http://gigo.gage.intranet/api/internal/ws/ext",
        method="POST",
        data=request_data.encode()
    )

    # send request to the GIGO server
    response = request.urlopen(req)

    # check response status
    if response.getcode() != 200:
        raise Exception(f"Failed to download extension with status code {response.getcode()}")

    # read response data and write to extension file
    with open("/tmp/gigo-developer.vsix", "w+b") as f:
        f.write(response.read())


def relay_failure(step: int, command: str, status: int, stdout: str, stderr: str, secret: str):
    """
    Uploads the failure status to the remote GIGO server for debugging
    :param step: step number from the execution instructions
    :param command: arbitrary shell command to be executed
    :param status: return status of the command
    :param stdout: stdout of the command
    :param stderr: stderr of the command
    :param secret: secret to be used for interacting with the remote GIGO server
    """
    # format request data
    request_data = json.dumps({
        "coder_id": args.workspace_id,
        "secret": secret,
        "step": step,
        "command": command,
        "status": status,
        "stdout": stdout,
        "stderr": stderr
    })

    # format POST request
    req = request.Request(
        "http://gigo.gage.intranet/api/internal/ws/init-failure",
        method="POST",
        data=request_data.encode()
    )

    # send request to the GIGO server
    request.urlopen(req)


def update_step(step: int, secret: str):
    """
    Updates the step number after a successful step completion via the remote GIGO server
    :param step: latest completed step number from the execution instructions
    :param secret: secret to be used for interacting with the remote GIGO server
    """
    # format request data
    request_data = json.dumps({
        "coder_id": args.workspace_id,
        "secret": secret,
        "step": step
    })

    # format POST request
    req = request.Request(
        "http://gigo.gage.intranet/api/internal/ws/init-step",
        method="POST",
        data=request_data.encode()
    )

    # send request to the GIGO server
    request.urlopen(req)


def has_been_initialized() -> bool:
    """
    has_been_initialized
    Checks if the workspace has been initialized.
    :return: bool
    """
    # check if this workspace has already been initialized
    return os.path.exists("/home/gigo/.gigo/ws-config.json")


def initialize_workspace() -> Tuple[dict, dict]:
    """
    initialize_workspace
    Initializes the workspace by calling the init function for the workspace
    in the GIGO internal API and returns the workspace config and git configuration.
    :return: Tuple[dict, dict] workspace_config, git_config
    """

    # create empty variable to hold data
    data: Optional[dict] = None

    try:
        # format request data
        request_data = json.dumps({"coder_id": args.workspace_id})

        # format POST request
        req = request.Request("http://gigo.gage.intranet/api/internal/ws/init", method="POST",
                              data=request_data.encode())

        # execute request
        res = request.urlopen(req)

        # read response and load json
        data = json.loads(res.read().decode())

        # extract git credentials
        git_config = {
            "git_email": data["git_email"],
            "git_token": data["git_token"],
            "git_name": data["owner_id"]
        }
        del data["git_email"]
        del data["git_token"]

        # update remote server for successful step
        update_step(1, data["secret"])
    except Exception as e:
        if data is not None and "secret" in data.keys():
            relay_failure(
                1,
                "remote initialization",
                -1,
                "",
                ''.join(tb.format_exception(None, e, e.__traceback__)),
                data["secret"]
            )
            exit(1)
        else:
            raise e

    return data, git_config


def write_gitconfig(git_config: dict, secret: str) -> None:
    """
    write_gitconfig
    Writes the git configuration to /home/gigo/.gitconfig.
    :param git_config: dict containing the git configuration
    :param secret: secret to be used for interacting with the remote GIGO server
    """
    try:
        # open .gitconfig and write configuration data
        with open(os.path.expanduser('/home/gigo/.gitconfig'), "w+") as f:
            f.write(f"""
[user]
	email = {git_config['git_email']}
	name = {git_config['git_name']}
[url "http://{git_config['git_name']}:{git_config['git_token']}@git.gage.intranet"]
	insteadOf = http://git.gage.intranet
"""
                    )

        # update remote server for successful step
        update_step(2, secret)
    except Exception as e:
        relay_failure(
            2,
            "write git config",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            secret
        )
        exit(1)


def write_workspace_config(workspace_config: dict) -> None:
    """
    write_workspace_config
    Writes the workspace configuration to /home/gigo/.gigo/ws-config.json
    :param workspace_config: dict containing the workspace configuration
    """
    try:
        # open workspace configuration file and save remaining data
        os.makedirs("/home/gigo/.gigo/", exist_ok=True)
        with open("/home/gigo/.gigo/ws-config.json", "w+") as f:
            json.dump(workspace_config, f)

        # update remote server for successful step
        update_step(3, workspace_config["secret"])
    except Exception as e:
        relay_failure(
            3,
            "write workspace config",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            workspace_config["secret"]
        )
        exit(1)


def clone_repo(workspace_config: dict):
    """
    clone_repo
    Clones the repository from the workspace configuration.
    :param workspace_config: dict containing the workspace configuration
    """
    # execute git clone command
    status, stdout, stderr = execute_command(
        f"git clone --recursive {workspace_config['repo']} {WORKING_DIRECTORY}")
    if status != 0:
        relay_failure(
            4,
            f"git clone --recursive {workspace_config['repo']} {WORKING_DIRECTORY}",
            status,
            stdout,
            stderr,
            workspace_config["secret"]
        )
        exit(1)

    try:
        # update remote server for successful step
        update_step(4, workspace_config["secret"])
    except Exception as e:
        relay_failure(
            4,
            f"git clone --recursive {workspace_config['repo']} {WORKING_DIRECTORY}",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            workspace_config["secret"]
        )
        exit(1)

    # execute checkout command
    status, stdout, stderr = execute_command(
        f"cd {WORKING_DIRECTORY} && git checkout {workspace_config['commit']}")
    if status != 0:
        relay_failure(
            5,
            f"cd {WORKING_DIRECTORY} && git checkout {workspace_config['commit']}",
            status,
            stdout,
            stderr,
            workspace_config["secret"]
        )
        exit(1)

    try:
        # update remote server for successful step
        update_step(5, workspace_config["secret"])
    except Exception as e:
        relay_failure(
            5,
            f"cd {WORKING_DIRECTORY} && git checkout {workspace_config['commit']}",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            workspace_config["secret"]
        )
        exit(1)


def handle_containers(secret: str, initialized: bool = False):
    """
    handle_containers
    Handles the containers in the workspace configuration by writing the compose file and bringing up the container set.
    :param secret: secret to be used for interacting with the remote GIGO server
    :param initialized: whether the workspace has been previously initialized
    """
    # skip container setup if the containers are not configured in this config
    if CONTAINER_COMPOSE.find("<containers>") != -1:
        return

    # initialize container directory
    try:
        os.makedirs("/home/gigo/.gigo/containers", exist_ok=True)
    except Exception as e:
        relay_failure(
            6,
            "create container directory",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            secret
        )
        exit(1)

    try:
        # update remote server for successful step
        update_step(6, secret)
    except Exception as e:
        relay_failure(
            6,
            "create container directory",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            secret
        )
        exit(1)

    # write container compose to container directory
    try:
        with open("/home/gigo/.gigo/containers/docker-compose.yml", "w+") as f:
            f.write(CONTAINER_COMPOSE)
    except Exception as e:
        relay_failure(
            7,
            "write container compose",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            secret
        )
        exit(1)

    try:
        # update remote server for successful step
        update_step(7, secret)
    except Exception as e:
        relay_failure(
            7,
            "write container compose",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            secret
        )
        exit(1)

    # execute container up command from container directory
    status, stdout, stderr = execute_command(
        f"cd /home/gigo/.gigo/containers && sudo docker-compose up -d"
    )
    if status != 0:
        relay_failure(
            8,
            f"cd /home/gigo/.gigo/containers && sudo docker-compose up -d",
            status,
            stdout,
            stderr,
            secret
        )
        exit(1)

    try:
        # update remote server for successful step
        update_step(8, secret)
    except Exception as e:
        relay_failure(
            8,
            f"cd /home/gigo/.gigo/containers && sudo docker-compose up -d",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            secret
        )
        exit(1)


def handle_user_executions(secret: str, initialized: bool = False):
    """
    handle_user_executions
    Handles the user executions in the workspace configuration by executing each command in the order specified.
    :param secret: secret to be used for interacting with the remote GIGO server
    :param initialized: whether the workspace has been previously initialized
    """
    # load shell executions from json string
    try:
        executions = json.loads(SHELL_EXECUTIONS)
    except Exception as e:
        relay_failure(
            9,
            "load shell executions",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            secret
        )
        exit(0)

    # iterate the commands executing each one
    for shell_execution in executions:
        # skip shell executions if it is only for initialization and we have already initialized
        if initialized and shell_execution["init"]:
            continue

        try:
            status, stdout, stderr = execute_command(shell_execution["command"])
            if status!= 0:
                relay_failure(
                    9,
                    shell_execution["name"],
                    status,
                    stdout,
                    stderr,
                    secret
                )
                exit(1)
        except Exception as e:
            relay_failure(
                9,
                shell_execution["name"],
                -1,
                "",
                ''.join(tb.format_exception(None, e, e.__traceback__)),
                secret
            )
            exit(1)

    try:
        # update remote server for successful step
        update_step(9, secret)
    except Exception as e:
        relay_failure(
            8,
            "update remote server",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            secret
        )
        exit(1)


def handle_vscode(_id: int, secret: str):
    """
    handle_vscode
    Handles the process of installing vscode (code-server), installing the configured extensions, and launching the execution
    :param _id: id of the workspace
    :param secret: secret to be used for interacting with the remote GIGO server
    :param initialized: whether the workspace has been previously initialized
    """
    # exit if we are skipping vscode
    if not USE_VSCODE:
        return

    # install vscode (code-server)
    try:
        status, stdout, stderr = execute_command("curl -fsSL https://code-server.dev/install.sh | sh")
        if status!= 0:
            relay_failure(
                10,
                "curl -fsSL https://code-server.dev/install.sh | sh",
                status,
                stdout,
                stderr,
                secret
            )
            exit(1)
        update_step(10, secret)
    except Exception as e:
        relay_failure(
            10,
            "install vscode",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            secret
        )
        exit(1)

    # iterate vscode extensions executing install
    for extension in VSCODE_EXTENSIONS:
        try:
            status, stdout, stderr = execute_command(f"code-server --install-extension {extension}")
            if status!= 0:
                relay_failure(
                    11,
                    f"code-server --install-extension {extension}",
                    status,
                    stdout,
                    stderr,
                    secret
                )
                exit(1)
        except Exception as e:
            relay_failure(
                11,
                f"failed to install vscode extension: {extension}",
                -1,
                "",
                ''.join(tb.format_exception(None, e, e.__traceback__)),
                secret
            )
            exit(1)

    # download gigo extension
    # install gigo developer extension
    try:
        download_extension(_id, secret)
    except Exception as e:
        relay_failure(
            11,
            f"failed to download gigo-dev extension",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            secret
        )
        exit(1)

    # install gigo developer extension
    try:
        status, stdout, stderr = execute_command(f"code-server --install-extension /tmp/gigo-developer.vsix")
        if status!= 0:
            relay_failure(
                11,
                f"code-server --install-extension /tmp/gigo-developer.vsix",
                status,
                stdout,
                stderr,
                secret
            )
            exit(1)
    except Exception as e:
        relay_failure(
            11,
            f"failed to install vscode extension: /tmp/gigo-developer.vsix",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            secret
        )
        exit(1)

    try:
        update_step(11, secret)
    except Exception as e:
        relay_failure(
            11,
            f"update remote server",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            secret
        )
        exit(1)

    try:
        update_step(12, secret)
    except Exception as e:
        relay_failure(
            12,
            f"update remote server",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            secret
        )
        exit(1)

    # launch vscode (code-server)
    try:
        status, stdout, stderr = execute_command("code-server --auth none --port 13337")
        if status!= 0:
            relay_failure(
                12,
                "code-server --auth none --port 13337",
                status,
                stdout,
                stderr,
                secret
            )
            exit(1)
    except Exception as e:
        relay_failure(
            12,
            "launch vscode",
            -1,
            "",
            ''.join(tb.format_exception(None, e, e.__traceback__)),
            secret
        )
        exit(1)


# execute script
if __name__ == "__main__":
    # check if remote initialization has occurred
    has_initialized = has_been_initialized()

    # conditionally execute workspace initialization steps
    if not has_initialized:
        ws_config, git_config = initialize_workspace()
        write_gitconfig(git_config, ws_config["secret"])
        write_workspace_config(ws_config)
        clone_repo(ws_config)
    else:
        # read existing workspace config
        try:
            with open("/home/gigo/.gigo/ws-config.json", "r") as f:
                ws_config = json.load(f)
        except Exception as e:
            relay_failure(
                14,
                "read workspace config",
                -1,
                "",
                ''.join(tb.format_exception(None, e, e.__traceback__)),
                ws_config["secret"]
            )
            exit(1)
        try:
            update_step(14, ws_config["secret"])
        except Exception as e:
            relay_failure(
                14,
                f"update remote server",
                -1,
                "",
                ''.join(tb.format_exception(None, e, e.__traceback__)),
                ws_config["secret"]
            )
            exit(1)


    # handle every start logic
    handle_containers(ws_config["secret"], has_initialized)
    handle_user_executions(ws_config["secret"], has_initialized)
    handle_vscode(ws_config["workspace_id_string"], ws_config["secret"])