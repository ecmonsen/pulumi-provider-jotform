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
    questions=[
        pulumi_jotform.FormQuestionsQuestionArgs(
            name="nameOfQ1",
            order=1,
            text="Enter your nickname.",
            type="control_textbox",
            properties={
                "hint": "Scooter"
            }
        ),
        pulumi_jotform.FormQuestionsQuestionArgs(
            name="nameOfQ2",
            order=2,
            text="Select a value.",
            type="control_radio",
            properties={
                "allowOther": "Yes",
                "options": "1|2|3|banana|5",
                "required": "true"
            }
        ),

    ],
    opts=pulumi.ResourceOptions(provider=provider))
