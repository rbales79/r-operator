# r-operator

A Go-based OpenShift Operator that configures Argo pipelines from custom resources.

## Usage

1. **Install the CRD:**

    ```sh
    kubectl apply -f config/crd/bases/pipelines.example.com_argopipelines.yaml
    ```

2. **Example ArgoPipeline resource:**

    ```yaml
    apiVersion: pipelines.example.com/v1
    kind: ArgoPipeline
    metadata:
      name: hello-pipeline
      namespace: default
    spec:
      pipelineYaml: |
        apiVersion: argoproj.io/v1alpha1
        kind: Workflow
        metadata:
          name: hello-world
        spec:
          entrypoint: whalesay
          templates:
          - name: whalesay
            container:
              image: docker/whalesay
              command: [cowsay]
              args: ["hello world"]
    ```

3. **Deploy the operator:**

    ```sh
    make run
    ```

4. **Result:** The operator will create/update the Argo Workflow as described in the `pipelineYaml` field.

## Requirements

- Go 1.19+
- Controller-runtime
- Operator SDK (for development)
- Argo Workflows installed on your cluster

## Notes

- The operator expects valid Argo Workflow YAML.
- The Argo CRDs must be installed in your cluster.
