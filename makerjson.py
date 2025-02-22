import json

with open("update.txt", "r", encoding="utf-8") as f:
    paths = [line.strip() for line in f if line.strip()]

with open("shell.json", "w", encoding="utf-8") as f:
    json.dump(paths, f, indent=2)

print("Convert Finish Saved At shell.json")
