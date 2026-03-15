package plain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	// DefaultBaseURL is the default URL for the Plain GraphQL API
	DefaultBaseURL = "https://core-api.uk.plain.com/graphql/v1"
	// DefaultTimeout is the default HTTP timeout
	DefaultTimeout = 30 * time.Second
)

// Client represents a Plain API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// NewClient creates a new Plain API client with the given auth token
func NewClient(token string) *Client {
	return &Client{
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		token: token,
	}
}

// NewClientWithURL creates a new client with a custom base URL (useful for testing)
func NewClientWithURL(token, baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		token: token,
	}
}

// Error represents an API error response
type Error struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
}

// request makes an HTTP request to the Plain API with authentication
func (c *Client) request(method, endpoint string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.baseURL
	if endpoint != "" {
		url = fmt.Sprintf("%s/%s", c.baseURL, endpoint)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		apiErr := &Error{
			StatusCode: resp.StatusCode,
			Message:    string(respBody),
		}

		// Try to parse error as JSON
		var errResp map[string]interface{}
		if json.Unmarshal(respBody, &errResp) == nil {
			if msg, ok := errResp["message"].(string); ok {
				apiErr.Message = msg
			} else if msg, ok := errResp["error"].(string); ok {
				apiErr.Message = msg
			}
		}

		return apiErr
	}

	// Parse successful response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

// Thread-related types

// Thread represents a support thread
type Thread struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	AssignedTo  *User      `json:"assignedTo,omitempty"`
	CreatedAt   DateTime   `json:"createdAt"`
	UpdatedAt   DateTime   `json:"updatedAt"`
	Labels      []Label    `json:"labels,omitempty"`
	Description string     `json:"description,omitempty"`
	Timeline    *Timeline  `json:"timeline,omitempty"`
}

// FormatPriority converts a priority integer to a human-readable string
func FormatPriority(priority int) string {
	switch priority {
	case 0:
		return "LOW"
	case 1:
		return "NORMAL"
	case 2:
		return "HIGH"
	case 3:
		return "URGENT"
	default:
		return "-"
	}
}

// DateTime represents a Plain API DateTime object
type DateTime struct {
	ISO8601 string `json:"iso8601"`
}

// Time converts DateTime to time.Time
func (dt DateTime) Time() (time.Time, error) {
	return time.Parse(time.RFC3339, dt.ISO8601)
}

// Label represents a thread label
type Label struct {
	ID        string     `json:"id"`
	LabelType *LabelType `json:"labelType,omitempty"`
}

// LabelType represents a label type
type LabelType struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Icon       string `json:"icon,omitempty"`
	Color      string `json:"color,omitempty"`
	IsArchived bool   `json:"isArchived"`
}

// User represents a Plain user
type User struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	FullName   string `json:"fullName"`
	PublicName string `json:"publicName,omitempty"`
}

// Attachment represents a file attachment
type Attachment struct {
	ID            string   `json:"id"`
	FileName      string   `json:"fileName"`
	FileSize      FileSize `json:"fileSize"`
	FileExtension string   `json:"fileExtension,omitempty"`
	FileMimeType  string   `json:"fileMimeType"`
	Type          string   `json:"type"`
	CreatedAt     DateTime `json:"createdAt"`
}

// FileSize represents file size information
type FileSize struct {
	Bytes int64 `json:"bytes"`
}

// AttachmentDownloadURL represents a temporary download URL for an attachment
type AttachmentDownloadURL struct {
	Attachment  Attachment `json:"attachment"`
	DownloadURL string     `json:"downloadUrl"`
	ExpiresAt   DateTime   `json:"expiresAt"`
}

// Actor represents an actor in timeline entries (union type)
type Actor struct {
	Typename string `json:"__typename"`
	// UserActor fields
	UserID   string `json:"userId,omitempty"`
	Email    string `json:"email,omitempty"`
	FullName string `json:"fullName,omitempty"`
	// CustomerActor fields
	CustomerID string `json:"customerId,omitempty"`
	// SystemActor fields
	SystemID string `json:"systemId,omitempty"`
	// MachineUserActor fields
	MachineUserID string `json:"machineUserId,omitempty"`
}

// ThreadFieldSchema represents a custom field schema for threads
type ThreadFieldSchema struct {
	ID                  string   `json:"id"`
	Key                 string   `json:"key"`
	Label               string   `json:"label"`
	Type                string   `json:"type"` // STRING, BOOL, ENUM, NUMBER, CURRENCY, DATE
	Description         string   `json:"description,omitempty"`
	EnumValues          []string `json:"enumValues,omitempty"`
	IsRequired          bool     `json:"isRequired"`
	IsAiAutoFillEnabled bool     `json:"isAiAutoFillEnabled"`
}

// Workspace represents a Plain workspace
type Workspace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// HelpCenter represents a help center
type HelpCenter struct {
	ID string `json:"id"`
}

// HelpCenterArticle represents a help center article
type HelpCenterArticle struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	ContentHTML string       `json:"contentHtml"`
	Slug        string       `json:"slug"`
	Status      string       `json:"status"`
	UpdatedAt   DateTime     `json:"updatedAt"`
	Group       *ArticleGroup `json:"articleGroup"`
}

// ArticleGroup represents an article group
type ArticleGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Timeline represents a thread timeline
type Timeline struct {
	Entries []TimelineEntry `json:"entries"`
}

// TimelineEntry represents a single timeline entry
type TimelineEntry struct {
	ID        string    `json:"id"`
	Timestamp DateTime  `json:"timestamp"`
	Actor     Actor     `json:"actor"`
	Entry     EntryData `json:"entry"`
}

// EntryData represents a timeline entry (union type)
type EntryData struct {
	Typename string `json:"__typename"`

	// NoteEntry fields
	NoteID      string       `json:"noteId,omitempty"`
	NoteText    string       `json:"noteText,omitempty"`
	Text        string       `json:"text,omitempty"` // fallback for compatibility
	Markdown    string       `json:"markdown,omitempty"`
	IsEdited    bool         `json:"isEdited,omitempty"`

	// ChatEntry fields
	ChatID         string    `json:"chatId,omitempty"`
	ChatText       string    `json:"chatText,omitempty"`
	CustomerReadAt *DateTime `json:"customerReadAt,omitempty"`

	// EmailEntry fields
	EmailID         string            `json:"emailId,omitempty"`
	Subject         string            `json:"subject,omitempty"`
	MarkdownContent string            `json:"markdownContent,omitempty"`
	From            *EmailParticipant `json:"from,omitempty"`
	To              *EmailParticipant `json:"to,omitempty"`

	// Status/Assignment change fields
	PreviousStatus   string          `json:"previousStatus,omitempty"`
	NextStatus       string          `json:"nextStatus,omitempty"`
	PreviousAssignee *ThreadAssignee `json:"previousAssignee,omitempty"`
	NextAssignee     *ThreadAssignee `json:"nextAssignee,omitempty"`
	PreviousPriority *int            `json:"previousPriority,omitempty"`
	NextPriority     *int            `json:"nextPriority,omitempty"`

	// SlackMessageEntry fields
	SlackText           string `json:"slackText,omitempty"`
	SlackWebMessageLink string `json:"slackWebMessageLink,omitempty"`

	// SlackReplyEntry fields
	SlackReplyText string `json:"slackReplyText,omitempty"`

	// ThreadDiscussionMessageEntry fields
	DiscussionText string `json:"discussionText,omitempty"`

	// ThreadDiscussionResolvedEntry fields
	ResolvedAt *DateTime `json:"resolvedAt,omitempty"`

	// ThreadDiscussionEntry fields
	DiscussionType string `json:"discussionType,omitempty"`

	// ThreadLabelsChangedEntry fields
	AddedLabelTypes   []LabelType `json:"addedLabelTypes,omitempty"`
	RemovedLabelTypes []LabelType `json:"removedLabelTypes,omitempty"`

	// ThreadAdditionalAssigneesTransitionedEntry fields
	NextAssignees     []ThreadAssignee `json:"nextAssignees,omitempty"`
	PreviousAssignees []ThreadAssignee `json:"previousAssignees,omitempty"`

	// ThreadLinkCreatedEntry fields
	Thread *struct {
		ID    string `json:"id,omitempty"`
		Title string `json:"title,omitempty"`
	} `json:"thread,omitempty"`

	// CustomEntry fields
	Title string `json:"title,omitempty"`
	Type  string `json:"type,omitempty"`

	// Common fields
	Attachments []Attachment `json:"attachments,omitempty"`
}

// EmailParticipant represents an email participant
type EmailParticipant struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// ThreadAssignee represents a thread assignee (union type)
type ThreadAssignee struct {
	Typename string `json:"__typename"`
	ID       string `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	FullName string `json:"fullName,omitempty"`
}

// GetEntryDisplayText returns a human-readable description of the entry
func (e *EntryData) GetEntryDisplayText() string {
	switch e.Typename {
	case "NoteEntry":
		return "Note"
	case "ChatEntry":
		return "Chat message"
	case "EmailEntry":
		return "Email"
	case "ThreadStatusTransitionedEntry":
		return "Status changed"
	case "ThreadAssignmentTransitionedEntry":
		return "Assignment changed"
	case "ThreadPriorityChangedEntry":
		return "Priority changed"
	case "SlackMessageEntry":
		return "Slack message"
	case "SlackReplyEntry":
		return "Slack reply"
	case "ThreadDiscussionMessageEntry":
		return "Discussion message"
	case "ThreadDiscussionResolvedEntry":
		return "Marked discussion as resolved"
	case "ThreadDiscussionEntry":
		return "Discussion started"
	case "ThreadLabelsChangedEntry":
		return "Labels changed"
	case "ThreadAdditionalAssigneesTransitionedEntry":
		return "Additional assignees changed"
	case "ThreadLinkCreatedEntry":
		return "Thread linked"
	default:
		return e.Typename
	}
}

// GetActorDisplayName returns the actor's display name
func (a *Actor) GetActorDisplayName() string {
	if a.FullName != "" {
		return a.FullName
	}
	if a.Email != "" {
		return a.Email
	}
	if a.Typename == "SystemActor" {
		return "System"
	}
	return "Unknown"
}

// Note represents an internal note on a thread
type Note struct {
	ID        string   `json:"id"`
	Text      string   `json:"text"`
	CreatedAt DateTime `json:"createdAt"`
	CreatedBy User     `json:"createdBy"`
}

// ThreadFilters contains filter options for listing threads
type ThreadFilters struct {
	Status     string   `json:"status,omitempty"`
	AssigneeID string   `json:"assigneeId,omitempty"`
	LabelIDs   []string `json:"labelIds,omitempty"`
	Priority   string   `json:"priority,omitempty"`
	Limit      int      `json:"limit,omitempty"`
	Offset     int      `json:"offset,omitempty"`

	// CreatedAfter filters threads created after this timestamp (ISO8601 format)
	// Example: "2026-03-15T00:00:00Z"
	CreatedAfter string `json:"-"`

	// CreatedBefore filters threads created before this timestamp (ISO8601 format)
	// Example: "2026-03-15T00:00:00Z"
	CreatedBefore string `json:"-"`

	// UpdatedAfter filters threads updated after this timestamp (ISO8601 format)
	// Example: "2026-03-15T00:00:00Z"
	UpdatedAfter string `json:"-"`

	// UpdatedBefore filters threads updated before this timestamp (ISO8601 format)
	// Example: "2026-03-15T00:00:00Z"
	UpdatedBefore string `json:"-"`
}

// ThreadsResponse represents a paginated list of threads
type ThreadsResponse struct {
	Threads []Thread `json:"threads"`
	Total   int      `json:"total"`
}

// graphql makes a GraphQL request to the Plain API
func (c *Client) graphql(query string, variables map[string]interface{}, result interface{}) error {
	body := map[string]interface{}{
		"query": query,
	}
	if variables != nil {
		body["variables"] = variables
	}
	return c.request("POST", "", body, result)
}

// Thread Operations

// ListThreads fetches a list of threads with optional filters
func (c *Client) ListThreads(filters *ThreadFilters) (*ThreadsResponse, error) {
	// Build the GraphQL query - based on working GetMyThreads query
	query := `
		query ListThreads($filters: ThreadsFilter, $first: Int, $after: String) {
			threads(filters: $filters, first: $first, after: $after) {
				edges {
					node {
						id
						title
						status
						priority
						assignedTo {
							... on User {
								id
								email
								fullName
							}
							... on MachineUser {
								id
								fullName
							}
						}
						labels {
							id
							labelType {
								id
								name
							}
						}
						createdAt {
							iso8601
						}
						updatedAt {
							iso8601
						}
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	// Build variables
	variables := map[string]interface{}{}

	// Build the filters object - matching GetMyThreads format
	threadFilters := map[string]interface{}{}

	if filters != nil {
		if filters.Status != "" {
			// Use "statuses" (plural) and convert to array
			threadFilters["statuses"] = []string{filters.Status}
		}

		if filters.AssigneeID != "" {
			// Use assignedToUser and convert to array
			threadFilters["assignedToUser"] = []string{filters.AssigneeID}
		}

		if filters.Priority != "" {
			// Use "priorities" (plural) - priority values must be integers
			// Convert string to int
			priority, err := strconv.Atoi(filters.Priority)
			if err != nil {
				return nil, fmt.Errorf("invalid priority value '%s': must be an integer", filters.Priority)
			}
			if priority > 0 {
				threadFilters["priorities"] = []int{priority}
			}
		}

		if len(filters.LabelIDs) > 0 {
			threadFilters["labelTypeIds"] = filters.LabelIDs
		}
	}

	if len(threadFilters) > 0 {
		variables["filters"] = threadFilters
	}

	// Set pagination
	limit := 50
	if filters != nil && filters.Limit > 0 {
		limit = filters.Limit
	}
	variables["first"] = limit

	// Note: GraphQL uses cursor-based pagination, not offset
	// For simplicity, we'll use offset as a workaround by using the cursor
	// In production, proper cursor handling should be implemented

	// Make the request - updated response struct to match query
	var response struct {
		Data struct {
			Threads struct {
				Edges []struct {
					Node struct {
						ID         string   `json:"id"`
						Title      string   `json:"title"`
						Status     string   `json:"status"`
						Priority   int      `json:"priority"`
						AssignedTo *User    `json:"assignedTo"`
						Labels     []Label  `json:"labels"`
						CreatedAt  DateTime `json:"createdAt"`
						UpdatedAt  DateTime `json:"updatedAt"`
					} `json:"node"`
				} `json:"edges"`
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"threads"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, variables, &response); err != nil {
		return nil, err
	}

	// Check for GraphQL errors
	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", response.Errors[0].Message)
	}

	// Convert response to ThreadsResponse
	threads := make([]Thread, 0, len(response.Data.Threads.Edges))
	for _, edge := range response.Data.Threads.Edges {
		thread := Thread{
			ID:         edge.Node.ID,
			Title:      edge.Node.Title,
			Status:     edge.Node.Status,
			Priority:   edge.Node.Priority,
			AssignedTo: edge.Node.AssignedTo,
			Labels:     edge.Node.Labels,
			CreatedAt:  edge.Node.CreatedAt,
			UpdatedAt:  edge.Node.UpdatedAt,
		}
		threads = append(threads, thread)
	}

	return &ThreadsResponse{
		Threads: threads,
		Total:   len(threads), // totalCount not available in pageInfo
	}, nil
}

// GetThread fetches details for a specific thread
func (c *Client) GetThread(threadID string, includeTimeline bool) (*Thread, error) {
	query := `
		query GetThread($threadId: ID!, $includeTimeline: Boolean!) {
			thread(threadId: $threadId) {
				id
				title
				description
				status
				priority
				assignedTo {
					... on User {
						id
						email
						fullName
					}
					... on MachineUser {
						id
						fullName
					}
				}
				labels {
					id
					labelType {
						id
						name
					}
				}
				createdAt {
					iso8601
				}
				updatedAt {
					iso8601
				}
				timelineEntries(first: 50) @include(if: $includeTimeline) {
					edges {
						node {
							id
							timestamp {
								iso8601
							}
							actor {
								__typename
								... on UserActor {
									userId
									user {
										id
										fullName
									}
								}
								... on CustomerActor {
									customerId
									customer {
										id
										fullName
									}
								}
								... on DeletedCustomerActor {
									customerId
								}
								... on SystemActor {
									systemId
								}
								... on MachineUserActor {
									machineUserId
									machineUser {
										id
										fullName
									}
								}
							}
							entry {
								__typename
								... on NoteEntry {
									noteId
									noteText: text
									markdown
									isEdited
									editedAt {
										iso8601
									}
									attachments {
										id
										fileName
										fileSize {
											bytes
										}
										fileExtension
										fileMimeType
										type
										createdAt {
											iso8601
										}
									}
								}
								... on ChatEntry {
									chatId
									chatText: text
									customerReadAt {
										iso8601
									}
									attachments {
										id
										fileName
										fileSize {
											bytes
										}
										fileExtension
										fileMimeType
										type
										createdAt {
											iso8601
										}
									}
								}
								... on EmailEntry {
									emailId
									subject
									textContent
									markdownContent
									from {
										name
										email
									}
									to {
										name
										email
									}
									additionalRecipients {
										name
										email
									}
									sentAt {
										iso8601
									}
									receivedAt {
										iso8601
									}
									attachments {
										id
										fileName
										fileSize {
											bytes
										}
										fileExtension
										fileMimeType
										type
										createdAt {
											iso8601
										}
									}
								}
								... on ThreadStatusTransitionedEntry {
									previousStatus
									nextStatus
								}
								... on ThreadAssignmentTransitionedEntry {
									previousAssignee {
										__typename
										... on User {
											id
											email
											fullName
										}
										... on MachineUser {
											id
											fullName
										}
									}
									nextAssignee {
										__typename
										... on User {
											id
											email
											fullName
										}
										... on MachineUser {
											id
											fullName
										}
									}
								}
								... on ThreadPriorityChangedEntry {
									previousPriority
									nextPriority
								}
								... on SlackMessageEntry {
									slackText: text
									slackWebMessageLink
									lastEditedOnSlackAt {
										iso8601
									}
									deletedOnSlackAt {
										iso8601
									}
									attachments {
										id
										fileName
										fileSize {
											bytes
										}
										fileExtension
										fileMimeType
										type
										createdAt {
											iso8601
										}
									}
								}
								... on CustomEntry {
									title
									type
									externalId
									attachments {
										id
										fileName
										fileSize {
											bytes
										}
										fileExtension
										fileMimeType
										type
										createdAt {
											iso8601
										}
									}
								}
								... on SlackReplyEntry {
									slackReplyText: text
									slackWebMessageLink
									attachments {
										id
										fileName
										fileSize {
											bytes
										}
										fileExtension
										fileMimeType
										type
										createdAt {
											iso8601
										}
									}
								}
								... on ThreadDiscussionMessageEntry {
									discussionText: text
								}
								... on ThreadDiscussionResolvedEntry {
									__typename
								}
								... on ThreadDiscussionEntry {
									__typename
								}
								... on ThreadLabelsChangedEntry {
									__typename
								}
								... on ThreadAdditionalAssigneesTransitionedEntry {
									nextAssignees {
										__typename
										... on User {
											id
											fullName
										}
										... on MachineUser {
											id
											fullName
										}
									}
									previousAssignees {
										__typename
										... on User {
											id
											fullName
										}
										... on MachineUser {
											id
											fullName
										}
									}
								}
								... on ThreadLinkCreatedEntry {
									__typename
								}
							}
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"threadId":        threadID,
		"includeTimeline": includeTimeline,
	}

	reqBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	var response struct {
		Data struct {
			Thread *struct {
				ID          string   `json:"id"`
				Title       string   `json:"title"`
				Description string   `json:"description"`
				Status      string   `json:"status"`
				Priority    int      `json:"priority"`
				AssignedTo  *User    `json:"assignedTo"`
				Labels      []Label  `json:"labels"`
				CreatedAt   DateTime `json:"createdAt"`
				UpdatedAt   DateTime `json:"updatedAt"`
				TimelineEntries *struct {
					Edges []struct {
						Node struct {
							ID        string    `json:"id"`
							Timestamp DateTime  `json:"timestamp"`
							Actor     struct {
								Typename string `json:"__typename"`
								// UserActor fields
								UserID string `json:"userId,omitempty"`
								User   *struct {
									ID       string `json:"id,omitempty"`
									FullName string `json:"fullName,omitempty"`
								} `json:"user,omitempty"`
								// CustomerActor fields
								CustomerID string `json:"customerId,omitempty"`
								Customer   *struct {
									ID       string `json:"id,omitempty"`
									FullName string `json:"fullName,omitempty"`
								} `json:"customer,omitempty"`
								// SystemActor fields
								SystemID string `json:"systemId,omitempty"`
								// MachineUserActor fields
								MachineUserID string `json:"machineUserId,omitempty"`
								MachineUser   *struct {
									ID       string `json:"id,omitempty"`
									FullName string `json:"fullName,omitempty"`
								} `json:"machineUser,omitempty"`
							} `json:"actor"`
							Entry EntryData `json:"entry"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"timelineEntries,omitempty"`
			} `json:"thread"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	err := c.request("POST", "", reqBody, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", response.Errors[0].Message)
	}

	if response.Data.Thread == nil {
		return nil, &Error{
			StatusCode: 404,
			Message:    fmt.Sprintf("Thread not found: %s", threadID),
		}
	}

	// Convert response to Thread
	thread := &Thread{
		ID:          response.Data.Thread.ID,
		Title:       response.Data.Thread.Title,
		Description: response.Data.Thread.Description,
		Status:      response.Data.Thread.Status,
		Priority:    response.Data.Thread.Priority,
		AssignedTo:  response.Data.Thread.AssignedTo,
		Labels:      response.Data.Thread.Labels,
		CreatedAt:   response.Data.Thread.CreatedAt,
		UpdatedAt:   response.Data.Thread.UpdatedAt,
	}

	// Convert timeline entries if they were fetched
	if response.Data.Thread.TimelineEntries != nil {
		entries := make([]TimelineEntry, 0, len(response.Data.Thread.TimelineEntries.Edges))
		for _, edge := range response.Data.Thread.TimelineEntries.Edges {
			// Flatten the actor structure
			actor := Actor{
				Typename: edge.Node.Actor.Typename,
			}

			// Map UserActor fields
			if edge.Node.Actor.User != nil {
				actor.UserID = edge.Node.Actor.UserID
				actor.FullName = edge.Node.Actor.User.FullName
			}

			// Map CustomerActor fields
			if edge.Node.Actor.Customer != nil {
				actor.CustomerID = edge.Node.Actor.CustomerID
				actor.FullName = edge.Node.Actor.Customer.FullName
			}

			// Map SystemActor fields
			if edge.Node.Actor.SystemID != "" {
				actor.SystemID = edge.Node.Actor.SystemID
			}

			// Map MachineUserActor fields
			if edge.Node.Actor.MachineUser != nil {
				actor.MachineUserID = edge.Node.Actor.MachineUserID
				actor.FullName = edge.Node.Actor.MachineUser.FullName
			}

			entries = append(entries, TimelineEntry{
				ID:        edge.Node.ID,
				Timestamp: edge.Node.Timestamp,
				Actor:     actor,
				Entry:     edge.Node.Entry,
			})
		}
		thread.Timeline = &Timeline{
			Entries: entries,
		}
	}

	return thread, nil
}

// SearchThreads searches for threads matching a query
func (c *Client) SearchThreads(query string, filters *ThreadFilters) (*ThreadsResponse, error) {
	// Build GraphQL query for searching threads
	graphqlQuery := `
		query SearchThreads($term: String!, $filters: ThreadsFilter, $first: Int) {
			searchThreads(
				searchQuery: { term: $term }
				filters: $filters
				first: $first
			) {
				edges {
					node {
						thread {
							id
							title
							status
							priority
							assignedTo {
								__typename
								... on User {
									id
									email
									fullName
								}
								... on MachineUser {
									id
									fullName
								}
							}
							labels {
								id
								labelType {
									id
									name
								}
							}
							createdAt {
								iso8601
							}
							updatedAt {
								iso8601
							}
						}
					}
				}
				pageInfo {
					hasNextPage
				}
			}
		}
	`

	// Build filters object
	filterMap := make(map[string]interface{})
	if filters != nil {
		if filters.Status != "" {
			filterMap["status"] = []string{filters.Status}
		}
		if filters.Priority != "" {
			filterMap["priority"] = []int{parsePriority(filters.Priority)}
		}
		if filters.AssigneeID != "" {
			filterMap["assignedToUser"] = []string{filters.AssigneeID}
		}
		if len(filters.LabelIDs) > 0 {
			filterMap["labelTypeIds"] = filters.LabelIDs
		}
	}

	// Build variables
	variables := map[string]interface{}{
		"term":  query,
		"first": 50, // Default limit
	}
	if len(filterMap) > 0 {
		variables["filters"] = filterMap
	}
	if filters != nil && filters.Limit > 0 {
		variables["first"] = filters.Limit
	}

	// Build request body
	requestBody := map[string]interface{}{
		"query":     graphqlQuery,
		"variables": variables,
	}

	// Make GraphQL request
	var result struct {
		Data struct {
			SearchThreads struct {
				Edges []struct {
					Node struct {
						Thread Thread `json:"thread"`
					} `json:"node"`
				} `json:"edges"`
				PageInfo struct {
					HasNextPage bool `json:"hasNextPage"`
				} `json:"pageInfo"`
			} `json:"searchThreads"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.request("POST", "", requestBody, &result); err != nil {
		return nil, err
	}

	// Check for GraphQL errors
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	// Convert to ThreadsResponse
	threads := make([]Thread, 0, len(result.Data.SearchThreads.Edges))
	for _, edge := range result.Data.SearchThreads.Edges {
		threads = append(threads, edge.Node.Thread)
	}

	return &ThreadsResponse{
		Threads: threads,
		Total:   len(threads),
	}, nil
}

// parsePriority converts a priority string to an integer value
func parsePriority(priority string) int {
	switch priority {
	case "urgent":
		return 0
	case "high":
		return 1
	case "normal":
		return 2
	case "low":
		return 3
	default:
		return 2 // Default to normal
	}
}

// Write Operations (Mutations)

// ChangeThreadStatus changes the workflow state of a thread (DONE, TODO, SNOOZED)
func (c *Client) ChangeThreadStatus(threadID string, status string) (*Thread, error) {
	mutation := `
		mutation ChangeThreadStatus($threadId: ID!, $status: ThreadStatus!) {
			changeThreadStatus(input: { threadId: $threadId, status: $status }) {
				thread {
					id
					title
					status
					priority
					updatedAt {
						iso8601
					}
				}
				error {
					message
					code
				}
			}
		}
	`

	variables := map[string]interface{}{
		"threadId": threadID,
		"status":   status,
	}

	var result struct {
		Data struct {
			ChangeThreadStatus struct {
				Thread Thread `json:"thread"`
				Error  *struct {
					Message string `json:"message"`
					Code    string `json:"code"`
				} `json:"error"`
			} `json:"changeThreadStatus"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(mutation, variables, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	if result.Data.ChangeThreadStatus.Error != nil {
		return nil, fmt.Errorf("%s", result.Data.ChangeThreadStatus.Error.Message)
	}

	return &result.Data.ChangeThreadStatus.Thread, nil
}

// SnoozeThread snoozes a thread until the specified datetime
func (c *Client) SnoozeThread(threadID string, until time.Time) (*Thread, error) {
	mutation := `
		mutation SnoozeThread($threadId: ID!, $until: DateTime!) {
			snoozeThread(input: { threadId: $threadId, until: $until }) {
				thread {
					id
					title
					status
					updatedAt {
						iso8601
					}
				}
				error {
					message
					code
				}
			}
		}
	`

	variables := map[string]interface{}{
		"threadId": threadID,
		"until":    until.Format(time.RFC3339),
	}

	var result struct {
		Data struct {
			SnoozeThread struct {
				Thread Thread `json:"thread"`
				Error  *struct {
					Message string `json:"message"`
					Code    string `json:"code"`
				} `json:"error"`
			} `json:"snoozeThread"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(mutation, variables, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	if result.Data.SnoozeThread.Error != nil {
		return nil, fmt.Errorf("%s", result.Data.SnoozeThread.Error.Message)
	}

	return &result.Data.SnoozeThread.Thread, nil
}

// AssignThread assigns a thread to a user (nil userID to unassign)
func (c *Client) AssignThread(threadID string, userID *string) (*Thread, error) {
	mutation := `
		mutation AssignThread($threadId: ID!, $userId: ID) {
			assignThreadToUser(input: { threadId: $threadId, userId: $userId }) {
				thread {
					id
					title
					assignedTo {
						__typename
						... on User {
							id
							email
							fullName
						}
						... on MachineUser {
							id
							fullName
						}
					}
					updatedAt {
						iso8601
					}
				}
				error {
					message
					code
				}
			}
		}
	`

	variables := map[string]interface{}{
		"threadId": threadID,
	}
	if userID != nil {
		variables["userId"] = *userID
	}

	var result struct {
		Data struct {
			AssignThreadToUser struct {
				Thread Thread `json:"thread"`
				Error  *struct {
					Message string `json:"message"`
					Code    string `json:"code"`
				} `json:"error"`
			} `json:"assignThreadToUser"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(mutation, variables, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	if result.Data.AssignThreadToUser.Error != nil {
		return nil, fmt.Errorf("%s", result.Data.AssignThreadToUser.Error.Message)
	}

	return &result.Data.AssignThreadToUser.Thread, nil
}

// ChangeThreadPriority changes the priority of a thread (0-3)
func (c *Client) ChangeThreadPriority(threadID string, priority int) (*Thread, error) {
	mutation := `
		mutation ChangeThreadPriority($threadId: ID!, $priority: Int!) {
			changeThreadPriority(input: { threadId: $threadId, priority: $priority }) {
				thread {
					id
					title
					priority
					updatedAt {
						iso8601
					}
				}
				error {
					message
					code
				}
			}
		}
	`

	variables := map[string]interface{}{
		"threadId": threadID,
		"priority": priority,
	}

	var result struct {
		Data struct {
			ChangeThreadPriority struct {
				Thread Thread `json:"thread"`
				Error  *struct {
					Message string `json:"message"`
					Code    string `json:"code"`
				} `json:"error"`
			} `json:"changeThreadPriority"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(mutation, variables, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	if result.Data.ChangeThreadPriority.Error != nil {
		return nil, fmt.Errorf("%s", result.Data.ChangeThreadPriority.Error.Message)
	}

	return &result.Data.ChangeThreadPriority.Thread, nil
}

// CreateNote creates an internal note on a thread
func (c *Client) CreateNote(threadID string, text string) (*Note, error) {
	mutation := `
		mutation CreateNote($threadId: ID!, $text: String!) {
			createNote(input: { threadId: $threadId, text: $text }) {
				note {
					id
					text
					createdAt {
						iso8601
					}
					createdBy {
						id
						email
						fullName
					}
				}
				error {
					message
					code
				}
			}
		}
	`

	variables := map[string]interface{}{
		"threadId": threadID,
		"text":     text,
	}

	var result struct {
		Data struct {
			CreateNote struct {
				Note  Note `json:"note"`
				Error *struct {
					Message string `json:"message"`
					Code    string `json:"code"`
				} `json:"error"`
			} `json:"createNote"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(mutation, variables, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	if result.Data.CreateNote.Error != nil {
		return nil, fmt.Errorf("%s", result.Data.CreateNote.Error.Message)
	}

	return &result.Data.CreateNote.Note, nil
}

// ReplyToThread sends a reply to a thread (customer-facing, not implemented in Phase 3)
func (c *Client) ReplyToThread(threadID, message string) error {
	return fmt.Errorf("not implemented: ReplyToThread (deferred to Phase 4)")
}

// ListLabelTypes fetches all available label types from the API
func (c *Client) ListLabelTypes(includeArchived bool) ([]*LabelType, error) {
	query := `
		query ListLabelTypes($isArchived: Boolean) {
			labelTypes(filters: { isArchived: $isArchived }) {
				edges {
					node {
						id
						name
						icon
						color
						isArchived
					}
				}
			}
		}
	`

	variables := map[string]interface{}{}
	if includeArchived {
		variables["isArchived"] = true
	}

	var result struct {
		Data struct {
			LabelTypes struct {
				Edges []struct {
					Node LabelType `json:"node"`
				} `json:"edges"`
			} `json:"labelTypes"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, variables, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	// Extract label types from edges
	labelTypes := make([]*LabelType, 0, len(result.Data.LabelTypes.Edges))
	for _, edge := range result.Data.LabelTypes.Edges {
		labelType := edge.Node
		labelTypes = append(labelTypes, &labelType)
	}

	return labelTypes, nil
}

// AddLabels adds label types to a thread
func (c *Client) AddLabels(threadID string, labelTypeIDs []string) ([]*Label, error) {
	mutation := `
		mutation AddLabels($threadId: ID!, $labelTypeIds: [ID!]!) {
			addLabels(input: { threadId: $threadId, labelTypeIds: $labelTypeIds }) {
				labels {
					id
					labelType {
						id
						name
						icon
						color
					}
				}
				error {
					message
					code
				}
			}
		}
	`

	variables := map[string]interface{}{
		"threadId":     threadID,
		"labelTypeIds": labelTypeIDs,
	}

	var result struct {
		Data struct {
			AddLabels struct {
				Labels []*Label `json:"labels"`
				Error  *struct {
					Message string `json:"message"`
					Code    string `json:"code"`
				} `json:"error"`
			} `json:"addLabels"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(mutation, variables, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	if result.Data.AddLabels.Error != nil {
		return nil, fmt.Errorf("%s", result.Data.AddLabels.Error.Message)
	}

	return result.Data.AddLabels.Labels, nil
}

// RemoveLabels removes specific label instances from a thread
func (c *Client) RemoveLabels(labelIDs []string) error {
	mutation := `
		mutation RemoveLabels($labelIds: [ID!]!) {
			removeLabels(input: { labelIds: $labelIds }) {
				thread {
					id
				}
				error {
					message
					code
				}
			}
		}
	`

	variables := map[string]interface{}{
		"labelIds": labelIDs,
	}

	var result struct {
		Data struct {
			RemoveLabels struct {
				Thread *Thread `json:"thread"`
				Error  *struct {
					Message string `json:"message"`
					Code    string `json:"code"`
				} `json:"error"`
			} `json:"removeLabels"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(mutation, variables, &result); err != nil {
		return err
	}

	if len(result.Errors) > 0 {
		return fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	if result.Data.RemoveLabels.Error != nil {
		return fmt.Errorf("%s", result.Data.RemoveLabels.Error.Message)
	}

	return nil
}

// ListThreadFieldSchemas fetches all available thread field schemas from the API
func (c *Client) ListThreadFieldSchemas() ([]*ThreadFieldSchema, error) {
	query := `
		query ListThreadFieldSchemas {
			threadFieldSchemas {
				edges {
					node {
						id
						key
						label
						type
						description
						enumValues
						isRequired
						isAiAutoFillEnabled
					}
				}
			}
		}
	`

	var result struct {
		Data struct {
			ThreadFieldSchemas struct {
				Edges []struct {
					Node ThreadFieldSchema `json:"node"`
				} `json:"edges"`
			} `json:"threadFieldSchemas"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, nil, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	// Extract field schemas from edges
	schemas := make([]*ThreadFieldSchema, 0, len(result.Data.ThreadFieldSchemas.Edges))
	for _, edge := range result.Data.ThreadFieldSchemas.Edges {
		schema := edge.Node
		schemas = append(schemas, &schema)
	}

	return schemas, nil
}

// GetToken returns the current auth token (useful for debugging)
func (c *Client) GetToken() string {
	return c.token
}

// SetTimeout sets the HTTP client timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// Help Center Operations

// ListHelpCenters fetches all available help centers
func (c *Client) ListHelpCenters() ([]*HelpCenter, error) {
	query := `
		query ListHelpCenters {
			helpCenters(first: 50) {
				edges {
					node {
						id
					}
				}
			}
		}
	`

	var result struct {
		Data struct {
			HelpCenters struct {
				Edges []struct {
					Node HelpCenter `json:"node"`
				} `json:"edges"`
			} `json:"helpCenters"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, nil, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	helpCenters := make([]*HelpCenter, 0, len(result.Data.HelpCenters.Edges))
	for _, edge := range result.Data.HelpCenters.Edges {
		hc := edge.Node
		helpCenters = append(helpCenters, &hc)
	}

	return helpCenters, nil
}

// ListWorkspaces fetches the user's workspace
func (c *Client) ListWorkspaces() ([]*Workspace, error) {
	query := `
		query GetMyWorkspace {
			myWorkspace {
				id
				name
			}
		}
	`

	var result struct {
		Data struct {
			MyWorkspace *Workspace `json:"myWorkspace"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, nil, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	// Check if myWorkspace is nil
	if result.Data.MyWorkspace == nil {
		return nil, fmt.Errorf("unable to fetch workspace: myWorkspace is nil")
	}

	// Return a slice with a single workspace
	return []*Workspace{result.Data.MyWorkspace}, nil
}

// GetHelpCenterArticle fetches a single article by ID
func (c *Client) GetHelpCenterArticle(articleID string) (*HelpCenterArticle, error) {
	query := `
		query HelpCenterArticle($id: ID!) {
			helpCenterArticle(id: $id) {
				id
				title
				contentHtml
				slug
				status
				updatedAt { iso8601 }
				articleGroup { id name }
			}
		}
	`

	variables := map[string]interface{}{
		"id": articleID,
	}

	var result struct {
		Data struct {
			Article *HelpCenterArticle `json:"helpCenterArticle"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, variables, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	if result.Data.Article == nil {
		return nil, &Error{
			StatusCode: 404,
			Message:    fmt.Sprintf("Article not found: %s", articleID),
		}
	}

	return result.Data.Article, nil
}

// ListHelpCenterArticles fetches all articles for a help center
func (c *Client) ListHelpCenterArticles(helpCenterID string, includeContent bool) ([]*HelpCenterArticle, error) {
	contentField := ""
	if includeContent {
		contentField = "contentHtml"
	}

	query := fmt.Sprintf(`
		query HelpCenterArticles($id: ID!, $first: Int, $after: String) {
			helpCenter(id: $id) {
				articles(first: $first, after: $after) {
					edges {
						node {
							id
							title
							slug
							status
							%s
							updatedAt { iso8601 }
							articleGroup { id name }
						}
					}
					pageInfo {
						hasNextPage
						endCursor
					}
				}
			}
		}
	`, contentField)

	var allArticles []*HelpCenterArticle
	cursor := ""

	for {
		variables := map[string]interface{}{
			"id":    helpCenterID,
			"first": 100,
		}
		if cursor != "" {
			variables["after"] = cursor
		}

		var result struct {
			Data struct {
				HelpCenter *struct {
					Articles struct {
						Edges []struct {
							Node HelpCenterArticle `json:"node"`
						} `json:"edges"`
						PageInfo struct {
							HasNextPage bool   `json:"hasNextPage"`
							EndCursor   string `json:"endCursor"`
						} `json:"pageInfo"`
					} `json:"articles"`
				} `json:"helpCenter"`
			} `json:"data"`
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		}

		if err := c.graphql(query, variables, &result); err != nil {
			return nil, err
		}

		if len(result.Errors) > 0 {
			return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
		}

		if result.Data.HelpCenter == nil {
			return nil, &Error{
				StatusCode: 404,
				Message:    fmt.Sprintf("Help center not found: %s", helpCenterID),
			}
		}

		for _, edge := range result.Data.HelpCenter.Articles.Edges {
			article := edge.Node
			allArticles = append(allArticles, &article)
		}

		if !result.Data.HelpCenter.Articles.PageInfo.HasNextPage {
			break
		}
		cursor = result.Data.HelpCenter.Articles.PageInfo.EndCursor
	}

	return allArticles, nil
}

// User Operations

// GetUserByEmail finds a user by their email address
func (c *Client) GetUserByEmail(email string) (*User, error) {
	query := `
		query GetUserByEmail($email: String!) {
			userByEmail(email: $email) {
				id
				email
				fullName
				publicName
			}
		}
	`

	variables := map[string]interface{}{
		"email": email,
	}

	var result struct {
		Data struct {
			UserByEmail *User `json:"userByEmail"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, variables, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	if result.Data.UserByEmail == nil {
		return nil, fmt.Errorf("user not found with email: %s", email)
	}

	return result.Data.UserByEmail, nil
}

// ListUsers fetches all users in the workspace
func (c *Client) ListUsers(assignableOnly bool) ([]*User, error) {
	query := `
		query ListUsers($filters: UsersFilter, $first: Int, $after: String) {
			users(filters: $filters, first: $first, after: $after) {
				edges {
					node {
						id
						email
						fullName
						publicName
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	var allUsers []*User
	cursor := ""

	for {
		variables := map[string]interface{}{
			"first": 100,
		}
		if assignableOnly {
			variables["filters"] = map[string]interface{}{
				"isAssignableToThread": true,
			}
		}
		if cursor != "" {
			variables["after"] = cursor
		}

		var result struct {
			Data struct {
				Users struct {
					Edges []struct {
						Node User `json:"node"`
					} `json:"edges"`
					PageInfo struct {
						HasNextPage bool   `json:"hasNextPage"`
						EndCursor   string `json:"endCursor"`
					} `json:"pageInfo"`
				} `json:"users"`
			} `json:"data"`
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		}

		if err := c.graphql(query, variables, &result); err != nil {
			return nil, err
		}

		if len(result.Errors) > 0 {
			return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
		}

		for _, edge := range result.Data.Users.Edges {
			user := edge.Node
			allUsers = append(allUsers, &user)
		}

		if !result.Data.Users.PageInfo.HasNextPage {
			break
		}
		cursor = result.Data.Users.PageInfo.EndCursor
	}

	return allUsers, nil
}

// CreateAttachmentDownloadUrl creates a temporary download URL for an attachment
func (c *Client) CreateAttachmentDownloadUrl(attachmentID string) (*AttachmentDownloadURL, error) {
	mutation := `
		mutation CreateAttachmentDownloadUrl($attachmentId: ID!) {
			createAttachmentDownloadUrl(input: { attachmentId: $attachmentId }) {
				attachmentDownloadUrl {
					attachment {
						id
						fileName
						fileSize {
							bytes
						}
						fileExtension
						fileMimeType
						type
						createdAt {
							iso8601
						}
					}
					downloadUrl
					expiresAt {
						iso8601
					}
				}
				error {
					message
					code
				}
			}
		}
	`

	variables := map[string]interface{}{
		"attachmentId": attachmentID,
	}

	var result struct {
		Data struct {
			CreateAttachmentDownloadUrl struct {
				AttachmentDownloadUrl *AttachmentDownloadURL `json:"attachmentDownloadUrl"`
				Error                 *struct {
					Message string `json:"message"`
					Code    string `json:"code"`
				} `json:"error"`
			} `json:"createAttachmentDownloadUrl"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(mutation, variables, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	if result.Data.CreateAttachmentDownloadUrl.Error != nil {
		return nil, fmt.Errorf("%s", result.Data.CreateAttachmentDownloadUrl.Error.Message)
	}

	if result.Data.CreateAttachmentDownloadUrl.AttachmentDownloadUrl == nil {
		return nil, &Error{
			StatusCode: 404,
			Message:    fmt.Sprintf("Attachment not found: %s", attachmentID),
		}
	}

	return result.Data.CreateAttachmentDownloadUrl.AttachmentDownloadUrl, nil
}

// DownloadAttachment downloads an attachment to the specified path
func (c *Client) DownloadAttachment(attachmentID, outputPath string) error {
	// Get download URL
	downloadInfo, err := c.CreateAttachmentDownloadUrl(attachmentID)
	if err != nil {
		return fmt.Errorf("failed to get download URL: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("GET", downloadInfo.DownloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() { _ = outFile.Close() }()

	// Copy response body to file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
