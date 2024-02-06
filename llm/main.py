import os
import uvicorn
from fastapi import FastAPI
from transformers import AutoTokenizer, AutoModel
import torch
from typing import List

# Create the app
app = FastAPI()

# Load the model name as an environment variable
model_name = os.getenv("MODEL_NAME", "bert-base-uncased")

# Load the tokeniser
tokenizer = AutoTokenizer.from_pretrained(model_name)

# Load the model
model = AutoModel.from_pretrained(model_name)

# Define a route to get the model configuration
@app.get("/config")
async def get_model_config():
    # Get the configuration of the model
    return model.config.__dict__

# Route for embedding text
@app.post("/embed")
async def embed_text(text_to_embed: List[str]) -> List[List[float]]:
    embeddings = []
    for text in text_to_embed:
        input_ids = tokenizer.encode(text, add_special_tokens=True)
        input_ids = torch.tensor(input_ids).unsqueeze(0)
        with torch.no_grad():
            output = model(input_ids)
        embeddings.append(output.last_hidden_state.mean(dim=1).numpy().flatten().tolist())
    print(embeddings)
    return embeddings

# Main loop
if __name__ == '__main__':
    uvicorn.run(app, host="0.0.0.0", port=8080)
