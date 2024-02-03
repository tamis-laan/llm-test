# Use the official PostgreSQL Alpine image as a base
FROM postgres:alpine
#
# Set the initial working directory to /tmp
WORKDIR /tmp

# Install the required tools
RUN apk --no-cache add git make curl gcc postgresql-dev musl-dev

# Clone the pg_vector repository (specific version v0.6.0)
RUN git clone --branch v0.6.0 https://github.com/pgvector/pgvector.git /tmp/pgvector
# Change into the pg_vector directory
WORKDIR /tmp/pgvector

# Build the extension
RUN make

# Install the extension (may require sudo depending on permissions)
RUN make install

# Clean up by removing the source code
WORKDIR /tmp
RUN rm -rf pgvector
