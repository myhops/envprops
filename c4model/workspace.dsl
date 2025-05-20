workspace {
    !identifiers hierarchical
    !impliedRelationships true

    model {
        modelViewer = person "Model Viewer" {
            tags InScope
        }

        c4modelSystem = softwareSystem "C4 Model System" {
            tags InScope
            structurizrOnPremises = container "Structurizr On Premises" {
                tags InScope
            }

            modelStorage = container "Model Storage" {
                tags InScope
            }

            modelCache = container "Model Cache" {
                tags InScope
            }

            sessionManager = container "Session Storage" {
                tags InScope
            }
            searchService = container "Search Service" {
                tags InScope
            }



            structurizrOnPremises -> modelStorage "Persist data using" 
            structurizrOnPremises -> modelCache "Cache models using" 
            structurizrOnPremises -> sessionManager "Share session data with" 
            structurizrOnPremises -> searchService "Index data using" 
        }

        modelViewer -> c4modelSystem.structurizrOnPremises "View Model"

        ssgIdentityProvider = softwareSystem "SSG Identity Provider" {
            samlService = container "SAML Service"
        }

        c4modelSystem.structurizrOnPremises -> ssgIdentityProvider.samlService "Authz Users at"
        

        productionEnv = deploymentEnvironment "Production" {
            aws = deploymentNode "AWS" {
                ec2 = deploymentNode "EC2 Instance" {
                    k0sCluster = deploymentNode "k0s Cluster" {
                        structurizrDeployment = deploymentNode "Structurizr Deployment" {
                            technology "Deployment"
                            structurizrPod = deploymentNode "Structurizr Pod" {
                                instances 1
                                c4modelInstance = containerInstance c4modelSystem.structurizrOnPremises
                            }
                        }
                    }
                }
                s3StorageNode = deploymentNode "S3 Storage Service" {
                    s3StoreageBucket = containerInstance c4modelSystem.modelStorage
                }
                elastiCacheNode = deploymentNode "ElastiCache Service" {
                    sessionCacheDB = containerInstance c4modelSystem.sessionManager {
                        description "Cache sessions"
                    }
                    modelCacheDB = containerInstance c4modelSystem.modelCache {
                        description "Cache sessions"
                    }
                }
                openSearchNode = deploymentNode "OpenSearch Service" {
                    elastiCacheDB = containerInstance c4modelSystem.searchService
                }
            }

            ssg = deploymentNode "SSG" {
                deploymentNode "Azure Cloud" {
                    deploymentNode "Entrada ID" {
                        samlService = deploymentNode "SSG SAML Service" {
                            samlService = containerInstance ssgIdentityProvider.samlService
                        }
                    }
                }
            }
        }

        ssgEnvironment = deploymentEnvironment "SSG" {
        }

        
    }

    views {
        
        systemContext c4modelSystem {
            include *
        }

        container c4modelSystem {
            include *
        }

        deployment * productionEnv {
            include * 
        }

        styles {
            element Person {
                shape Person
            }

            element InScope {
                background lightblue
            }
        }
    }
}

