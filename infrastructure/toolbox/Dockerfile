FROM golang:1.19

ENV USE_GKE_GCLOUD_AUTH_PLUGIN=True

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl \
    && chmod +x ./kubectl \
    && mv ./kubectl /usr/local/bin/kubectl

ARG CLOUD_SDK_VERSION=412.0.0
ENV PATH /google-cloud-sdk/bin:$PATH

RUN apt-get clean \
      && apt-get update -qq  \
      && apt-get install -y apt-utils debian-archive-keyring build-essential libpq-dev vim curl make

RUN curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz \
      && tar -xzf google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz -C / \
      && rm google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz \
      && /google-cloud-sdk/bin/gcloud components install gke-gcloud-auth-plugin \
      && /google-cloud-sdk/bin/gcloud config set core/disable_usage_reporting true \
      && /google-cloud-sdk/bin/gcloud config set component_manager/disable_update_check true \
      && /google-cloud-sdk/bin/gcloud config set metrics/environment github_docker_image

RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 \
      && chmod +x get_helm.sh \
      && ./get_helm.sh

RUN curl -s https://fluxcd.io/install.sh | bash

RUN mkdir /app
WORKDIR /app
