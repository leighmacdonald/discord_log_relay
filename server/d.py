import os

for f in os.listdir("."):
    if f.contains("z"):
        os.remove(f)