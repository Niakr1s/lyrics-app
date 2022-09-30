import argparse
import pathlib
from mutagen.id3 import ID3, USLT
from mutagen.oggvorbis import OggVorbis


def setLyrics(musicFileName: str, lyrics: str):
    ext = pathlib.Path(musicFileName).suffix
    if ext == ".mp3":
        setLyricsMp3(musicFileName, lyrics)
    elif ext == ".ogg":
        setLyricsOgg(musicFileName, lyrics)
    else:
        raise Exception("wrong extension, want .mp3 or .ogg")


def setLyricsMp3(musicFileName: str, lyrics: str):
    audio = ID3(musicFileName)
    audio.add(USLT(lang='   ', desc='', text=lyrics))
    audio.save()


def setLyricsOgg(musicFileName: str, lyrics: str):
    audio = OggVorbis(musicFileName)
    audio['lyrics'] = lyrics
    audio.save()


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Set lyrics to file.')
    parser.add_argument(
        "musicFilePath", help="Input music filepath. Must be .mp3 or .ogg", type=str)
    parser.add_argument("lyrics", help="lyrics string", type=str)
    args = parser.parse_args()

    musicFilePath = args.musicFilePath
    lyrics = args.lyrics

    setLyrics(musicFilePath, lyrics)

    print(musicFilePath, lyrics)
