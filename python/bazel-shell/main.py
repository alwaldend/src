import os 
import sys
import itertools
import shlex

if __name__ == "__main__":
    if working_dir := os.getenv("BAZEL_WORKING_DIRECTORY"):
        os.chdir(working_dir)
    os.execv("/bin/bash", ("/bin/bash", "-c", shlex.join(sys.argv[1:])))
