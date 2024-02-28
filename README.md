# mail-operator

## Description
A Kubernetes Operator for automating Email delivery.

## Getting Started

### Run the operator on a Kubernetes Cluster

To run the operator you need to make sure you have a cluster configured on your ~/.kube/config file.

Having that, run the following command to install the Email and EmailSenderConfig CRDs into your Kubernetes cluster:
```sh
make manifests; make generate; make install
```

Update the config/secrets/token_secret.yaml file to include your Mailersend and/or Mailergun API keys hashed as B64. Please keep their name the same as it will be referenced inside the application. You can use the following command to generate the hash:
```sh
echo -n "YOUR-API-KEY" | base64
```

You can run `kubectl apply -k config/default` to apply the Kustomize charts if you wish to use the [image](https://hub.docker.com/repository/docker/jbiers/controller/general) I deployed on Dockerhub. Otherwise, run the following commands to generate your own version of the image, and be sure to also update the image value in the config/manager/manager.yaml file.
```sh
docker build . -t <YOUR-REPOSITORY>/controller:latest
docker push <YOUR-REPOSITORY>/controller:latest
```

Now, the operator should be running on your cluster. For testing it, first create an EmailSenderConfig resource as per the example in config/samples/mail_v1_emailsenderconfig.yaml. Note that the value for *apiToken* should be either *mailersend* or *mailgun*, depending on the providers you have configured and which one you want to use. If you wish to use both for different emails, create two separate EmailSenderConfig resources. Create the resource with the command:
```sh
kubectl apply -f config/samples/mail_v1_emailsenderconfig.yaml
```

The remaining step is to create the Email resource itself. Edit the file at config/samples/mail_v1_email.yaml according to your scenario and create it.
```sh
kubectl apply -f config/samples/mail_v1_email.yaml
```

The email should now have been sent.

## About the project design

This project was built using the [Kubebuilder SDK](https://book.kubebuilder.io/). It is considered a [best practice](https://cloud.google.com/blog/products/containers-kubernetes/best-practices-for-building-kubernetes-operators-and-stateful-apps) in Operator development to use an SDK among the available options. Kubebuilder was chosen as I was more familiar with it.

The basic directory structure is generated by the SDK, and it goes as follows:

- *api/*      - contains the Golang files defining the CRDs for the project, in this case Email and EmailSenderConfig.
- *bin/*      - useful binaries for the project
- *cmd/*      - where the *main.go* file resides, being the entrypoint for the application
- *config/*   - contains all of the Kustomize charts for building the necessary resources for the project.
- *hack/*     - apache license files
- *internal/* - the actual code for the operator controller

I chose to develop the controller in a way that is easy to add new providers. You can find Mailersend and Mailgun as the two alternatives I implemented as an example, but a new one could be added by defining a struct and a sendEmail method for it (as per the interface defined in internal/provider/provider.go), having it as an option in the switch statement (internal/controller/email_controller.go, line 75) and adding its API key as a secret.

## Improvement points
This project was built in only 5 days. With more time in my hands, the following points should be addressed:

- Unit and Integration testing.
- Better test the MailGun integration. I did not have the chance to test it as extensively as Mailersend.
- Automate the deployment with a CI/CD tool like Github Actions.
- Generate Helm Charts. Kustomize does a good job here but Helm is my option of choice in most production cases.
