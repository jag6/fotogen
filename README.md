# Fotogen

## Description
This is a photo sharing web app that allows signed-in users to create galleries and upload images either from their device or Dropbox. Images are uploaded concurrently, and only one of each image is allowed per gallery. Deletion of gallery also deletes its images from disk.

### Feature and Functionality Overview
- Galleries
  - can only view if signed-in 
  - CRUD
  - Images
    - Multi-file upload from either device or oauth2 through Dropbox
    - duplicate images not allowed
- Reset password
- Search
  - query users by username
  - single input in navbar
  - display not found message if exact match not found, otherwise show link to user gallery index page
