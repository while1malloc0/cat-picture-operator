# Cat Picture Operator
We'll be creating a new custom resource type called a CatPicture. This will be acted upon by an operator written with the operator framework.

## Requirements

* A Go development environment
* A working Kubernetes cluster
* THe operator-sdk installed


## Build steps

### Create the project scaffolding

```bash
operator-sdk new operator --api-version=cat.example.com/v1 --kind=CatPicture --skip-git-init
```

