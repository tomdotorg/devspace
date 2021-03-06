---
title: 2. Containerize your app
---

DevSpace.cli lets you easily containerize your application, so you can deploy it to Kubernetes. You can use one of your own projects for the remainder of this guide. 

If you do not have a project to work with, you can **checkout one of our demo projects (optional)**
<!--DOCUSAURUS_CODE_TABS-->
<!--Node.js-->
```bash
git clone https://github.com/devspace-cloud/devspace-quickstart-nodejs
cd devspace-quickstart-nodejs
```

<!--END_DOCUSAURUS_CODE_TABS-->

> You can also use any existing project. DevSpace works with any programming language.

## Initialize your project
Run the following command within your project:
```bash
devspace init
```
DevSpace.cli will automatically detect your programming language and ask for the ports your application is listening on. Then, it will automatically create the following files:
```bash
project/                    # your project directory
|
|--.devspace/               # DevSpace directory
|   |-config.yaml           # DevSpace config
|
|--chart/                   # Helm chart (defines how to deploy your application)
|   |-Chart.yaml            # chart definition (e.g. name, version)
|   |-values.yaml           # values for the template variables
|   |-templates/            # directory containing the template files
|
|--Dockerfile               # Dockerfile (defines how to build the Docker image)
```

<details>
<summary>
### Learn how to customize Helm chart and image building (optional)
</summary>

See the following guides to:
- [Configure image building](../deployment/images)
- [Add packages to your Helm chart (e.g. database)](../charts/packages)
- [Configure persistent volumes](../charts/persistent-volumes)
- [Set environment variables](../charts/environment-variables)
- [Configure networking for your Helm chart (e.g. ingress)](../charts/networking)
- [Define multiple containers in your Helm chart](../charts/containers)
- [Add custom Kubernetes manifests (.yaml files)](../charts/custom-manifests)
- [Configure auto-scaling within your Helm Chart](../charts/scaling)

</details>
