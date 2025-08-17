# Requirements Document

## Introduction

The Junk Journal & Board web application is a private-by-URL digital journaling platform that allows users to create visual journal entries with drag-and-drop canvas editing capabilities. The system provides a simple, authentication-free approach where boards are accessed via secret tokens in URLs, making it easy to share and collaborate without complex user management. Users can create boards containing multiple pages, each with various elements like text, images, stickers, and shapes, all editable through an intuitive canvas interface with theme customization and export capabilities.

## Requirements

### Requirement 1: Board Management

**User Story:** As a user, I want to create and manage boards so that I can organize my journal entries into separate collections.

#### Acceptance Criteria

1. WHEN a user creates a new board THEN the system SHALL generate a unique edit token and public token
2. WHEN a board is created THEN the system SHALL return both edit URL and public URL for sharing
3. WHEN a user accesses a board with edit token THEN the system SHALL allow full editing capabilities
4. WHEN a user accesses a board with public token THEN the system SHALL provide read-only access
5. WHEN a user updates board metadata THEN the system SHALL require a valid edit token
6. WHEN a board is accessed THEN the system SHALL load all associated pages and elements

### Requirement 2: Page Management

**User Story:** As a user, I want to create and organize pages within boards so that I can structure my journal entries chronologically or thematically.

#### Acceptance Criteria

1. WHEN a user creates a page THEN the system SHALL require a title, date, and order index
2. WHEN pages are displayed THEN the system SHALL order them by the order_idx field
3. WHEN a user updates a page THEN the system SHALL require a valid edit token
4. WHEN a page is deleted THEN the system SHALL cascade delete all associated elements
5. WHEN a page is accessed THEN the system SHALL load all its elements with proper z-ordering

### Requirement 3: Canvas Element Management

**User Story:** As a user, I want to add and manipulate various elements on my journal pages so that I can create rich, visual content.

#### Acceptance Criteria

1. WHEN a user adds an element THEN the system SHALL support text, image, sticker, and shape types
2. WHEN an element is created THEN the system SHALL store position (x,y), dimensions (w,h), rotation, and z-index
3. WHEN elements are manipulated THEN the system SHALL support drag, resize, rotate, and reorder operations
4. WHEN element changes occur THEN the system SHALL debounce API calls by 300ms to optimize performance
5. WHEN multiple elements are selected THEN the system SHALL allow batch operations
6. WHEN elements are reordered THEN the system SHALL provide batch z-index update endpoint

### Requirement 4: Canvas Editing Features

**User Story:** As a user, I want advanced canvas editing capabilities so that I can precisely control the layout and appearance of my journal elements.

#### Acceptance Criteria

1. WHEN editing text elements THEN the system SHALL provide font family, size, color, bold, and italic controls
2. WHEN snap-to-grid is enabled THEN the system SHALL align elements to 8px grid increments
3. WHEN undo/redo is triggered THEN the system SHALL maintain a client-side history stack of maximum 50 operations
4. WHEN Delete key is pressed THEN the system SHALL remove selected elements
5. WHEN elements are selected THEN the system SHALL show visual selection indicators and manipulation handles

### Requirement 5: Theme and Customization

**User Story:** As a user, I want to apply different themes to my boards so that I can customize the visual appearance to match my preferences.

#### Acceptance Criteria

1. WHEN a user selects a theme THEN the system SHALL apply predefined background skins like wood or notebook
2. WHEN theme is changed THEN the system SHALL persist the selection to the board
3. WHEN themes are displayed THEN the system SHALL load options from the public/skins directory
4. WHEN a theme is applied THEN the system SHALL update the board's skin field

### Requirement 6: File Upload and Management

**User Story:** As a user, I want to upload images to use in my journal pages so that I can include personal photos and graphics.

#### Acceptance Criteria

1. WHEN a user uploads an image THEN the system SHALL validate file type (jpg, png, gif) and size (max 10MB)
2. WHEN an image is uploaded THEN the system SHALL store it in organized directories by board ID
3. WHEN upload is successful THEN the system SHALL return a public URL for the image
4. WHEN images are stored THEN the system SHALL use UUID-based filenames to prevent conflicts
5. WHEN upload fails THEN the system SHALL return appropriate error messages

### Requirement 7: Recap and Overview

**User Story:** As a user, I want to view summaries of my journal activity so that I can track my journaling habits and quickly find specific entries.

#### Acceptance Criteria

1. WHEN accessing recap view THEN the system SHALL provide day, week, and month filter options
2. WHEN a time range is selected THEN the system SHALL show page count and element count for that period
3. WHEN recap is displayed THEN the system SHALL show page thumbnails and metadata
4. WHEN filtering by date THEN the system SHALL use the page's date field for accurate results
5. WHEN recap data is requested THEN the system SHALL return structured JSON with counts and page details

### Requirement 8: Export Functionality

**User Story:** As a user, I want to export my journal pages as PNG images so that I can save and share my creations outside the application.

#### Acceptance Criteria

1. WHEN a user exports a page THEN the system SHALL generate a PNG image client-side using html2canvas
2. WHEN export is triggered THEN the system SHALL render at 2x resolution for high quality output
3. WHEN export completes THEN the system SHALL provide download functionality
4. WHEN exporting THEN the system SHALL capture the current state of all visible elements
5. WHEN export fails THEN the system SHALL display appropriate error messages

### Requirement 9: URL-Based Access Control

**User Story:** As a user, I want secure access to my boards without traditional authentication so that I can easily share while maintaining privacy.

#### Acceptance Criteria

1. WHEN accessing edit functionality THEN the system SHALL require edit_token in URL parameters
2. WHEN accessing public view THEN the system SHALL use public_token for read-only access
3. WHEN tokens are invalid THEN the system SHALL return appropriate error responses
4. WHEN tokens are generated THEN the system SHALL ensure uniqueness across all boards
5. WHEN API mutations are attempted THEN the system SHALL validate edit_token before processing

### Requirement 10: Data Persistence and API

**User Story:** As a user, I want my journal data to be reliably saved and accessible so that I don't lose my work and can access it from different sessions.

#### Acceptance Criteria

1. WHEN data is modified THEN the system SHALL persist changes to PostgreSQL database
2. WHEN API requests are made THEN the system SHALL return consistent JSON error format
3. WHEN database operations occur THEN the system SHALL use proper indexing for performance
4. WHEN concurrent access happens THEN the system SHALL handle race conditions appropriately
5. WHEN system starts THEN the system SHALL run database migrations automatically