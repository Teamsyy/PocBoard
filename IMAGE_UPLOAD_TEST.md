# Image Upload Testing Instructions

## Test 17: Image Upload Functionality

### Features to Test:

1. **File Upload via Button Click**
   - Click the "Upload Image" button in toolbar or sidebar
   - Select JPG, PNG, or GIF files
   - Verify file validation (max 10MB, correct types)
   - Check upload progress indicator
   - Confirm automatic image element creation on canvas

2. **Drag and Drop Upload**
   - Drag image files from computer over the browser
   - Should show blue drop zone overlay
   - Drop files to upload
   - Verify multiple file support
   - Check error handling for invalid files

3. **Clipboard Paste Upload** 
   - Copy an image from any application (e.g., screenshot, copied image)
   - Focus on the board editor page
   - Press Ctrl+V (or Cmd+V on Mac)
   - Should automatically upload and create image element
   - Verify pasted images are given auto-generated filenames

4. **Persistence Testing**
   - Upload/paste images
   - Refresh the page or navigate away and back
   - Verify uploaded images are still displayed
   - Check that image URLs are properly stored

5. **Error Handling**
   - Try uploading files over 10MB
   - Try uploading non-image files
   - Test with corrupted image files
   - Verify error messages are displayed

6. **UI/UX Testing**
   - Check button styling matches toolbar theme
   - Verify upload progress shows correctly
   - Test success/error toasts appear and auto-dismiss
   - Confirm drag overlay appears only when dragging files

### Expected Behavior:
- ✅ Images upload successfully via all three methods
- ✅ Images automatically appear on canvas at reasonable size
- ✅ Upload progress is shown during file transfer
- ✅ Error messages for validation failures
- ✅ Success confirmation messages
- ✅ Uploaded images persist across page reloads
- ✅ Clipboard paste works for screenshots and copied images
- ✅ Global drag/drop works anywhere on the page

### Backend Requirements:
- Upload endpoint: POST `/api/v1/boards/{boardId}/upload`
- File validation: JPG, PNG, GIF up to 10MB
- Organized storage: `/uploads/boards/{boardId}/`
- Static file serving for uploaded images

### Current Implementation Status: ✅ COMPLETE
All features implemented and integrated with existing canvas editor.
