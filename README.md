A plugin to download folders and files to Jfrog artifactory.

Run the following script to install git-leaks support to this repo.
```
chmod +x ./git-hooks/install.sh
./git-hooks/install.sh
```

# Building

Build the plugin binary:

```text
scripts/build.sh
```

Build the plugin image:

```text
docker build -t plugins/artifactory-download  -f docker/Dockerfile .
```

# Testing

Execute the plugin from your current working directory:

```text
docker run --rm \
  -e PLUGIN_USERNAME='username' \
  -e PLUGIN_PASSWORD='pwd' \
  -e PLUGIN_URL='artifactory instance url' \
  -e PLUGIN_SOURCE_PATH=Venkat-Test-BT-274/ \
  -e PLUGIN_TARGET_PATH=./harness/dbops/ \
  -e PLUGIN_INCLUDE_DIRS=true \
  -v $(pwd):/drone \
  plugins/artifactory-download:latest
```

## Harness CI Example:
```yaml
              - step:
                  type: Plugin
                  name: jFrog-Test
                  identifier: Artifactory_Download_Plugin
                  spec:
                    connectorRef: account.harnessImage
                    image: plugins/artifactory-download:linux-amd64
                    settings:
                      username: username
                      password: <JFROG_PWD>
                      url: https://URL.jfrog.io/artifactory
                      source_path: /harness/cache.txt
                      target_path: newdemo/
```

## Community and Support
[Harness Community Slack](https://join.slack.com/t/harnesscommunity/shared_invite/zt-y4hdqh7p-RVuEQyIl5Hcx4Ck8VCvzBw) - Join the #drone slack channel to connect with our engineers and other users running Drone CI.

[Harness Community Forum](https://community.harness.io/) - Ask questions, find answers, and help other users.

[Report and Track A Bug](https://community.harness.io/c/bugs/17) - Find a bug? Please report in our forum under Drone Bugs. Please provide screenshots and steps to reproduce. 

[Events](https://www.meetup.com/harness/) - Keep up to date with Drone events and check out previous events [here](https://www.youtube.com/watch?v=Oq34ImUGcHA&list=PLXsYHFsLmqf3zwelQDAKoVNmLeqcVsD9o).
