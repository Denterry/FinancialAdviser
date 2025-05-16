#!/usr/bin/env python3
import os
import subprocess
from pathlib import Path

def generate_grpc_code():
    """Generate gRPC code from proto files."""
    # Get the project root directory
    project_root = Path(__file__).parent.parent
    
    # Define paths
    proto_dir = project_root / "src" / "proto"
    proto_file = proto_dir / "ml_service.proto"
    
    # Create proto directory if it doesn't exist
    proto_dir.mkdir(parents=True, exist_ok=True)
    
    # Generate Python code
    subprocess.run([
        "python", "-m", "grpc_tools.protoc",
        f"--proto_path={proto_dir}",
        f"--python_out={proto_dir}",
        f"--grpc_python_out={proto_dir}",
        str(proto_file),
    ])
    
    print("Generated gRPC code successfully!")


if __name__ == "__main__":
    generate_grpc_code() 