import taglib

with taglib.File("test/Ayano Mashiro - Gentou.mp3", save_on_exit=True) as song:
    song.tags
    song.length
    song.tags["ALBUM"] = ["White Album"]
    del song.tags["DATE"]
    song.tags["GENRE"] = ["Vocal", "Classical"]
    song.tags["PERFORMER:HARPSICHORD"] = ["Ton Koopman"]
