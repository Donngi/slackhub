package tool

import (
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

// Tool is Struct which reflects tool's metadata in DynamoDB
type Tool struct {
	ID              string        `dynamo:"id"`
	DisplayName     string        `dynamo:"displayName"`
	Description     string        `dynamo:"description"`
	ModalJSON       string        `dynamo:"modalJSON"`
	CalleeArn       string        `dynamo:"calleeArn"`
	Administrators  []string      `dynamo:"administrators"`
	AuthorizedUsers []string      `dynamo:"authorizedUsers"`
	BootMode        string        `dynamo:"bootMode"`
	dbClient        toolInterface `dynamo:"-"`
}

// New is a constractor of Tool.
func New() *Tool {
	t := new(Tool)
	t.dbClient = &dbClient{}
	return t
}

// NewTool is a constractor of Tool.
func NewTool(id string, displayName string, description string, modalJSON string, calleeArn string, administrators []string, authorizedUsers []string, bootMode string) *Tool {
	t := Tool{
		ID:              id,
		DisplayName:     displayName,
		Description:     description,
		ModalJSON:       modalJSON,
		CalleeArn:       calleeArn,
		Administrators:  administrators,
		AuthorizedUsers: authorizedUsers,
		BootMode:        bootMode,
		dbClient:        &dbClient{},
	}
	return &t
}

// GetItem gets tool's metadata from DynamoDB and set them to struct.
func (t *Tool) GetItem(id string, region string, dynamodb string) error {

	var db = dynamo.New(session.New(), &aws.Config{
		Region: aws.String(region),
	})
	var table = db.Table(dynamodb)

	// Get selected tool's metadata from DynamoDB
	err := t.dbClient.getOne(table, "id", id, t)
	if err != nil {
		return err
	}

	// Set dbClient
	t.dbClient = &dbClient{}

	return nil
}

// GetAllItems gets all tool's metadata from DynamoDB.
func GetAllItems(region string, dynamodb string) ([]Tool, error) {

	var db = dynamo.New(session.New(), &aws.Config{
		Region: aws.String(region),
	})
	var table = db.Table(dynamodb)

	// Get all tool's metadata.
	var toolList []Tool
	err := scanAll(table, &toolList)
	if err != nil {
		return nil, err
	}

	// Set dbClient.
	var result []Tool
	for _, t := range toolList {
		t.dbClient = &dbClient{}
		result = append(result, t)
	}

	return result, nil
}

// IsAdministrators checks if a user is permitted to edit the tool.
func (t *Tool) IsAdministrators(user string) bool {
	if len(t.Administrators) == 0 {
		return true
	}

	for _, v := range t.Administrators {
		if user == v {
			return true
		}
	}
	return false
}

// IsAuthorizedUsers checks if a user is permitted to use the tool.
func (t *Tool) IsAuthorizedUsers(user string) bool {
	if len(t.AuthorizedUsers) == 0 {
		return true
	}

	for _, v := range t.AuthorizedUsers {
		if user == v {
			return true
		}
	}
	return false
}

// IsUseModal check wether the tool uses a modal or not.
func (t *Tool) IsUseModal() bool {
	if len(t.ModalJSON) == 0 {
		return false
	}
	return true
}

// Register puts new tool's metadata to DynamoDB.
func (t *Tool) Register(region string, dynamodb string) error {

	var db = dynamo.New(session.New(), &aws.Config{
		Region: aws.String(region),
	})
	var table = db.Table(dynamodb)

	// Register new item only when there isn't same id tool.
	if err := t.dbClient.putIf(table, t, "attribute_not_exists(id)"); err != nil {
		return err
	}

	return nil
}

// RegisterForce puts new tool's metadata to DynamoDB even if there is a same id tool.
func (t *Tool) RegisterForce(region string, dynamodb string) error {

	var db = dynamo.New(session.New(), &aws.Config{
		Region: aws.String(region),
	})
	var table = db.Table(dynamodb)

	// Register new item even if there ia a same id tool.
	if err := t.dbClient.put(table, t); err != nil {
		return err
	}

	return nil
}

// Delete delete an existing item.
func (t *Tool) Delete(region string, dynamodb string) error {
	var db = dynamo.New(session.New(), &aws.Config{
		Region: aws.String(region),
	})
	var table = db.Table(dynamodb)

	// Delete item.
	if err := t.dbClient.delete(table, "id", t.ID, t); err != nil {
		return err
	}
	return nil
}

// removeDuplication returns a slice which contains only no-duplicated value.
func removeDuplication(slice []string) []string {

	m := make(map[string]bool)
	var res []string

	for _, item := range slice {
		if _, ok := m[item]; !ok {
			m[item] = true
			res = append(res, item)
		}
	}

	return res
}

// removeElement returns a slice which eliminates the target element from original slice.
func removeElement(slice []string, target string) []string {

	var res []string

	for _, item := range slice {
		if item != target {
			res = append(res, item)
		}
	}
	return res
}

// removeElements returns a slice which eliminates the target element from original slice.
func removeElements(slice []string, targets []string) []string {

	res := slice

	for _, t := range targets {
		res = removeElement(res, t)
	}
	return res
}

// SortTools returns the sorted slice of Tool. SlackHub's official tools are located at the end.
func SortTools(toolList []Tool) []Tool {
	// To avoid inserting empty struct, use some slices
	var res []Tool
	var register, editor, catalog, eraser []Tool
	for _, v := range toolList {
		switch v.ID {
		case "register":
			register = append(register, v)
		case "editor":
			editor = append(editor, v)
		case "catalog":
			catalog = append(catalog, v)
		case "eraser":
			eraser = append(eraser, v)
		default:
			res = append(res, v)
		}
	}

	// Sort
	sort.Slice(res, func(i, j int) bool {
		return res[i].DisplayName < res[j].DisplayName
	})

	res = append(res, register...)
	res = append(res, editor...)
	res = append(res, catalog...)
	res = append(res, eraser...)

	return res
}
