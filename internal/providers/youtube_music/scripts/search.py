import sys
import json
from ytmusicapi import YTMusic

def main():
    query = sys.argv[1]
    ytmusic = YTMusic()  # unauthenticated is fine for search
    results = ytmusic.search(query, filter="songs", limit=5)

    simplified = [
        {
            "title": r["title"],
            "artist": r["artists"][0]["name"] if r.get("artists") else "",
            "videoId": r["videoId"],
            "duration": r.get("duration", ""),
        }
        for r in results
    ]
    print(json.dumps(simplified))

if __name__ == "__main__":
    main()
    