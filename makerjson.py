import json

# Baca daftar path dari file
with open("update.txt", "r", encoding="utf-8") as f:
    paths = [line.strip() for line in f if line.strip()]

# Simpan dalam format JSON
with open("paths.json", "w", encoding="utf-8") as f:
    json.dump(paths, f, indent=2)

print("Konversi selesai! File tersimpan sebagai paths.json")
