#!/bin/bash
echo "Running code generation with wire..."
cd cmd/api && wire
echo "Wire generation completed."
