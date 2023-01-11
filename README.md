# stuff to do

## api endpoints

- /api/slideshow/next - retrieve next photo metadata to show
- /api/upload - upload a new photo?
- /api/list - get all photos

## frontend routes

- /slideshow - load the TV slideshow
- / - welcome screen - option to view photos or upload
- /upload - place to upload photo
- /browse - photo browser

## server capabilities

- scan raw photos and find ones that are not yet in the slideshow
  - use timestamp of last copied?
  - all photos newer than last recorded timestamp are copied in
  - then update timestamp
- process newly found photos
  - if shot in RAW - convert to jpeg
  - resize for full screen and thumbnail
  - add to photo database by filename
- photo database
  - records photo url, last shown timestamp (0 if new photo)
  - uses mysqlite for small database

## frontend capabilities

- slideshow photo transitions
- slideshow alert for new photo
- background color for TV
- browse option - show list of thumbnails and allow for full screen view
- pagination?
- welcome screen
- upload screen
