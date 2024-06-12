"""Deployment Pipeline."""
from aws_cdk import App, Environment
from gaggle_cdk.core import GaggleTags

from stacks.pipeline_stack import PipelineStack

app = App()

toolchain_account = app.node.try_get_context("gaggle-cdk:toolchain-account")
toolchain_region = app.node.try_get_context("gaggle-cdk:toolchain-region")

PipelineStack(
    app,
    name="couture-simulator",
    environment="toolchain",
    team=GaggleTags.Team.TOOL,
    env=Environment(account=toolchain_account, region=toolchain_region),
)

app.synth()
