import argparse
import pathlib
from mutagen.id3 import ID3, USLT
from mutagen.oggvorbis import OggVorbis


def writeLyrics(musicFileName: str, lyrics: str):
    ext = pathlib.Path(musicFileName).suffix
    if ext == ".mp3":
        writeLyricsMp3(musicFileName, lyrics)
    elif ext == ".ogg":
        writeLyricsOgg(musicFileName, lyrics)
    else:
        raise Exception("wrong extension, want .mp3 or .ogg")


def writeLyricsMp3(musicFileName: str, lyrics: str):
    audio = ID3(musicFileName)
    audio.add(USLT(lang='   ', desc='', text=lyrics))
    audio.save()


def writeLyricsOgg(musicFileName: str, lyrics: str):
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

    writeLyrics(musicFilePath, lyrics)

    print("done write lyrics to", musicFilePath)
