#!/bin/bash
#SBATCH --time=30:00:00
#SBATCH --nodes=1
#SBATCH --cpus-per-task=5
#SBATCH --job-name=WICSIE
#SBATCH --mem=20GB

module load Go

# Create directories
mkdir out
mkdir out/raw

sh cleanOut.sh

# Compile
go run main.go