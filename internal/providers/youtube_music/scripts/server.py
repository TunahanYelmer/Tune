import sys
import json
from ytmusicapi import YTMusic

def main():
    # Initialize once — this is the expensive part we're avoiding repeating.
    ytmusic = YTMusic()

    # Signal readiness so Go knows it's safe to start sending requests.
    print(json.dumps({"ready": True}), flush=True)

    for line in sys.stdin:
        line = line.strip()
        if not line:
            continue

        try:
            req = json.loads(line)
            query = req.get("query", "")
            limit = req.get("limit", 5)

            results = ytmusic.search(query, filter="songs", limit=limit)
            results = results[:limit]

            simplified = [
                {
                    "title": r.get("title", ""),
                    "artist": r["artists"][0]["name"] if r.get("artists") else "",
                    "album": r["album"]["name"] if r.get("album") else "",
                    "duration": r.get("duration", ""),
                    "videoId": r.get("videoId", ""),
                }
                for r in results
                if r.get("videoId")
            ]
            print(json.dumps({"tracks": simplified}), flush=True)

        except Exception as e:
            print(json.dumps({"error": str(e)}), flush=True)

if __name__ == "__main__":
    main()