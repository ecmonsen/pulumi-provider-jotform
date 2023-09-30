"""A Python Pulumi program"""

import pulumi
import pulumi_jotform

config = pulumi.Config()

jotform_api_key = config.require("jotform_api_key")
provider = pulumi_jotform.Provider("myprovider", api_key=jotform_api_key)

myform = pulumi_jotform.Form(
    "myform",
    title="My Form",
    emails=[
        pulumi_jotform.EmailArgsArgs(type="foo", from_="abc@def", to="ghi@jkl", subject="testing", html=True,
                                     body="<p>hi</p>")
    ],
    opts=pulumi.ResourceOptions(provider=provider))
pulumi.export("myform.formid", myform.form_id)

myform_questions = pulumi_jotform.FormQuestions(
    "myform-questions",
    form_id=myform.form_id,
    questions=[],
    opts=pulumi.ResourceOptions(provider=provider))
