#!/usr/bin/env python3
import os
import sys
import subprocess
import time

def print_header(message):
    """Print a formatted header message."""
    print("\n" + "=" * 80)
    print(f" {message} ".center(80, "="))
    print("=" * 80 + "\n")

def run_script(script_path, args=None):
    """Run a Python script with optional arguments."""
    cmd = [sys.executable, script_path]
    if args:
        cmd.extend(args)
    
    print(f"Running: {' '.join(cmd)}")
    
    try:
        start_time = time.time()
        result = subprocess.run(cmd, check=True)
        elapsed_time = time.time() - start_time
        print(f"Completed in {elapsed_time:.2f} seconds with exit code {result.returncode}")
        return True
    except subprocess.CalledProcessError as e:
        print(f"Error: Script failed with exit code {e.returncode}")
        return False
    except Exception as e:
        print(f"Error: {str(e)}")
        return False

def main():
    """Run all maintenance scripts."""
    base_dir = os.path.dirname(os.path.abspath(__file__))
    
    # List of scripts to run in order
    scripts = [
        "ensure_directory_structure.py",
        "fix_latex_formatting.py",
        "fix_local_links.py",
    ]
    
    print_header("STARTING MAINTENANCE PROCESS")
    
    success_count = 0
    failure_count = 0
    
    for script in scripts:
        script_path = os.path.join(base_dir, script)
        if not os.path.exists(script_path):
            print(f"Warning: Script {script} not found at {script_path}")
            continue
        
        print_header(f"RUNNING {script}")
        
        if run_script(script_path, [base_dir]):
            success_count += 1
        else:
            failure_count += 1
    
    print_header("MAINTENANCE COMPLETE")
    print(f"Scripts completed successfully: {success_count}")
    print(f"Scripts failed: {failure_count}")
    
    if failure_count > 0:
        print("\nSome scripts failed. Please check the output above for details.")
        return 1
    else:
        print("\nAll maintenance tasks completed successfully!")
        return 0

if __name__ == "__main__":
    sys.exit(main()) 