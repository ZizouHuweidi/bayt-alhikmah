import os
from llama_index.llms.openai import OpenAI
from llama_index.core import PromptTemplate


OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")
if not OPENAI_API_KEY:
    raise ValueError("Please set the OPENAI_API_KEY environment variable.")

# Initialize the OpenAI LLM
llm = OpenAI(model="gpt-4o-mini", temperature=0.7)

# Define the drug recommendation prompt template
RECOMMENDATION_PROMPT = """

"""

# Create a PromptTemplate
prompt_template = PromptTemplate(RECOMMENDATION_PROMPT)


# Function to generate recommendations
def get_recommendation(user_input):
    # Format the prompt with the user input
    formatted_prompt = prompt_template.format(user_input=user_input)

    # Generate the response using OpenAI
    response = llm.complete(formatted_prompt)

    return response.text


# Test the function in a Jupyter Notebook
user_input = ""
recommendation = get_recommendation(user_input)
print(recommendation)
