"""Pipeline Stack."""
from aws_cdk import Stack, Aws
from constructs import Construct
# from gaggle_cdk.core.pipelines.pipeline_utils import ecr_permissions
# from gaggle_cdk.core.pipelines.pipeline_utils import codeartifact_permissions,
from gaggle_cdk.core import GaggleTags
from gaggle_cdk.core.pipelines import DeploymentPipeline, DeploymentEnvironment


class PipelineStack(Stack):
    def __init__(
            self,
            scope: Construct,
            name: str,
            environment: str,
            team: GaggleTags.Team,
            **kwargs,
    ) -> None:
        super().__init__(scope, f"{name}-PipelineStack", **kwargs)

        DeploymentPipeline(
            self,
            application_name=name,
            environment=environment,
            github_connection_param="/toolchain/github/connection_arn",
            github_repo="cdk-app-template",  # Change to your Github Repo
            github_org="gaggle-net",
            main_branch="main",
            team=team,
            # ecr_repo_name=["gaggle/app-repo"], # Change to your ECR repo name(s)
            slack_channel_config_param="/toolchain/chatbot/slack/<team-name>/arn",  # Change to your team name
            build_vars={
                # Add Environment variables for your Build stages here
                # Secrets Manager secrets can be added like: "MY_SECRET": "secret:secret/id"
                # Parameter Store parameters can be added like: `"MY_PARAM": "parameter:parameter/name"
            },
            integration_environment=DeploymentEnvironment(
                account_id="123456789012",  # Replace with Integration Account ID
                hosted_zone="integration-foo.gaggle.services",  # Replace with Integration Account domain name
                env_vars={
                    "CDK_ACCOUNT": "123456789012",  # Replace with Integration Account ID
                    "CDK_REGION": Aws.REGION,  # Replace with Integration Account Region if not same
                    "ENVIRONMENT": "integration-foo",  # Replace with Integration Environment name
                    "DEPLOY_IT_ROLE": "true",
                },
            ),
            staging_environment=DeploymentEnvironment(
                account_id="123456789012",  # Replace with Staging Account ID
                hosted_zone="staging.gaggle.services",  # Replace with Staging Account domain name
                env_vars={
                    "CDK_ACCOUNT": "123456789012",  # Replace with Staging Account ID
                    "CDK_REGION": Aws.REGION,  # Replace with Staging Account Region if not same
                    "ENVIRONMENT": "staging",  # Replace with Staging Environment name
                },
            ),
            production_environment=DeploymentEnvironment(
                account_id="123456789012",  # Replace with Production Account ID
                env_vars={
                    "CDK_ACCOUNT": "123456789012",  # Replace with Production Account ID
                    "CDK_REGION": Aws.REGION,  # Replace with Production Account Region if not same
                    "ENVIRONMENT": "production",  # Replace with Production Environment name
                },
            ),
        )

        # If you need ECR permissions, uncomment this
        # ecr_permissions(
        #     codebuild_project=pipeline.ci_build,
        #     ecr_account=core.Aws.ACCOUNT_ID,
        #     ecr_region=core.Aws.REGION,
        #     publish=True,
        # )

        # ecr_permissions(
        #     codebuild_project=pipeline.build_project,
        #     ecr_account=core.Aws.ACCOUNT_ID,
        #     ecr_region=core.Aws.REGION,
        #     publish=True,
        # )

        # If you need CodeArtifact permissions, uncomment this
        # codeartifact_permissions(
        #     codebuild_project=pipeline.ci_build,
        #     codeartifact_account=core.Aws.ACCOUNT_ID,
        #     codeartifact_region=core.Aws.REGION,
        #     publish=True,
        # )

        # codeartifact_permissions(
        #     codebuild_project=pipeline.build_project,
        #     codeartifact_account=core.Aws.ACCOUNT_ID,
        #     codeartifact_region=core.Aws.REGION,
        #     publish=True,
        # )
