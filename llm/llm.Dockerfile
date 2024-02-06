# Use the official Python image as the base image
FROM python:slim

# Set environment variables for FastAPI
ENV APP_HOME /app

# Set the working directory
WORKDIR $APP_HOME

# Update apt
RUN apt update

# Install packages
RUN apt install -y fd-find entr make curl

# Install specific version of MarkupSafe (compatible with jinja2)
RUN pip install --no-cache-dir MarkupSafe==2.0

# Install pytorch
RUN pip install torch --index-url https://download.pytorch.org/whl/cpu

# Install huggingface transformers
RUN pip install transformers

# Install fastapi
RUN pip install uvicorn fastapi

# Set the model
ENV MODEL_NAME roberta-base

# Download the model and tokenizer during the Docker image build
RUN python -c "from transformers import AutoTokenizer, AutoModel; \
	tokenizer = AutoTokenizer.from_pretrained('$MODEL_NAME'); \
	model = AutoModel.from_pretrained('$MODEL_NAME')"

# Copy the FastAPI application code into the container
COPY . .

# Expose port 8000
EXPOSE 8080

# Command to run the FastAPI application using pipenv
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]
