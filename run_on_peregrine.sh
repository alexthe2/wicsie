#!/bin/bash
#SBATCH --time=00:10:00
#SBATCH --nodes=1
#SBATCH --cpus-per-task=4
#SBATCH --job-name=WICSIE_test
#SBATCH --mem=10GB

module load Go/1.16.6

# Create directories
mkdir out
mkdir out/raw

# Compile
go run main.go