# %% Cell 1,
from dotenv import dotenv_values
import os
import requests
import pandas as pd

from langchain.chains import RetrievalQA
from langchain.document_loaders import TextLoader
from langchain.embeddings.openai import OpenAIEmbeddings
from langchain.llms import OpenAI,HuggingFaceHub
from langchain.text_splitter import CharacterTextSplitter
from langchain.vectorstores import Chroma, FAISS
from langchain.document_loaders.csv_loader import CSVLoader
from langchain.prompts import PromptTemplate
from langchain.embeddings import HuggingFaceEmbeddings

config = dotenv_values(".env")
HUGGINGFACEHUB_API_TOKEN = config['HUGGINGFACEHUB_API_TOKEN']
os.environ["HUGGINGFACEHUB_API_TOKEN"] = HUGGINGFACEHUB_API_TOKEN

# %% Cell 2,

df = pd.read_csv('./data/goodreads_data.csv').drop(['Unnamed: 0'],axis=1)
df = df.assign( genre_len = lambda x:len(x['Genres']))
df = df[df['genre_len']>0]
df = df[['Book','Genres','Description']]
df.to_csv('book_genre.csv',index=False)
print(df.head())
