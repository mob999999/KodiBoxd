import xbmc
import xbmcgui
import xbmcaddon
import json
import urllib.parse
import urllib.request
import os

ADDON = xbmcaddon.Addon()
ADDON_NAME = ADDON.getAddonInfo('name')
ADDON_PATH = xbmc.translatePath(ADDON.getAddonInfo('path'))
QR_CODE_PATH = os.path.join(ADDON_PATH, 'qr_code.png')

def get_library_data():
    """Fetch the user's Kodi library data."""
    json_query = '{"jsonrpc": "2.0", "method": "VideoLibrary.GetMovies", "params": {"properties": ["title"]}, "id": 1}'
    response = xbmc.executeJSONRPC(json_query)
    data = json.loads(response)
    
    json_query_tv = '{"jsonrpc": "2.0", "method": "VideoLibrary.GetTVShows", "params": {"properties": ["title"]}, "id": 1}'
    response_tv = xbmc.executeJSONRPC(json_query_tv)
    data_tv = json.loads(response_tv)
    
    movies = [movie['title'] for movie in data.get('result', {}).get('movies', [])]
    tv_shows = [show['title'] for show in data_tv.get('result', {}).get('tvshows', [])]
    
    return {"movies": movies, "tv_shows": tv_shows}

def generate_qr_code():
    """Generate a QR code using goqr.me API and display it in Kodi."""
    library_data = get_library_data()
    json_data = json.dumps(library_data)
    encoded_json = urllib.parse.quote(json_data)
    
    qr_url = f"https://api.qrserver.com/v1/create-qr-code/?data=https://google.com/{encoded_json}&size=300x300"
    
    try:
        urllib.request.urlretrieve(qr_url, QR_CODE_PATH)
        xbmcgui.Dialog().ok(ADDON_NAME, "QR code generated. Now displaying.")
        xbmc.executebuiltin(f"ShowPicture({QR_CODE_PATH})")
    except Exception as e:
        xbmcgui.Dialog().ok(ADDON_NAME, f"Error generating QR code: {str(e)}")

if __name__ == '__main__':
    generate_qr_code()
