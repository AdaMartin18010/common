#!/usr/bin/env python3
import os
import subprocess
import sys
import time

def run_script(script_path, description):
    """Run a script and capture its output."""
    print(f"\n{'='*80}")
    print(f"Running: {description}")
    print(f"{'='*80}\n")
    
    start_time = time.time()
    
    try:
        result = subprocess.run(['python', script_path], 
                                capture_output=True, 
                                text=True)
        
        print(result.stdout)
        
        if result.stderr:
            print("Errors:")
            print(result.stderr)
        
        elapsed = time.time() - start_time
        print(f"\nFinished in {elapsed:.2f} seconds with exit code: {result.returncode}")
        
        return result.returncode == 0
    except Exception as e:
        print(f"Failed to run {script_path}: {e}")
        return False

def main():
    """Run all maintenance scripts in sequence."""
    refactor_dir = "docs/refactor"
    
    if not os.path.isdir(refactor_dir):
        print(f"Error: {refactor_dir} directory not found.")
        sys.exit(1)
    
    scripts = [
        ("docs/refactor/ensure_directory_structure.py", "Checking directory structure"),
        ("docs/refactor/verify_links.py", "Verifying internal links"),
        ("docs/refactor/fix_all_math_expressions.py", "Fixing LaTeX math expressions")
    ]
    
    results = []
    
    for script_path, description in scripts:
        success = run_script(script_path, description)
        results.append((description, success))
    
    # Print summary
    print("\n\n")
    print(f"{'='*80}")
    print("MAINTENANCE SUMMARY")
    print(f"{'='*80}")
    
    all_success = True
    for description, success in results:
        status = "✅ SUCCESS" if success else "❌ FAILED"
        print(f"{status} - {description}")
        all_success = all_success and success
    
    if all_success:
        print("\n✅ All maintenance tasks completed successfully!")
        return 0
    else:
        print("\n❌ Some maintenance tasks failed. Please check the output above.")
        return 1

if __name__ == "__main__":
    sys.exit(main()) 