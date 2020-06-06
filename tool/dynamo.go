package tool

// All methods in this file is lapper of DynamoDB access.
// It makes unit test easy.

import "github.com/guregu/dynamo"

type toolInterface interface {
	getOne(table dynamo.Table, key string, value string, t *Tool) error
	put(table dynamo.Table, t *Tool) error
	putIf(table dynamo.Table, t *Tool, condition string) error
	delete(table dynamo.Table, key string, value string, t *Tool) error
}

type dbClient struct{}

// getOne is a lapper method of DynamoDB access. Get a tool's metadata from DynamoDB.
func (c *dbClient) getOne(table dynamo.Table, key string, value string, t *Tool) error {
	if err := table.Get(key, value).One(t); err != nil {
		return err
	}
	return nil
}

// scanAll is a lapper method of DynamoDB access. Get all tool's metadata from DynamoDB.
var scanAll = func(table dynamo.Table, toolList *[]Tool) error {
	if err := table.Scan().All(toolList); err != nil {
		return err
	}
	return nil
}

// put is a lapper method of DynamoDB access. Put a tool's metadata to DynamoDB.
func (c *dbClient) put(table dynamo.Table, t *Tool) error {
	if err := table.Put(t).Run(); err != nil {
		return err
	}
	return nil
}

// putIf is a lapper method of DynamoDB access. Put a tool's metadata to DynamoDB only in specific conditions.
func (c *dbClient) putIf(table dynamo.Table, t *Tool, condition string) error {
	if err := table.Put(t).If(condition).Run(); err != nil {
		return err
	}
	return nil
}

// delete is a lapper method of DynamoDB access. Put a tool's metadata to DynamoDB only in specific conditions.
func (c *dbClient) delete(table dynamo.Table, key string, value string, t *Tool) error {
	if err := table.Delete(key, value).Run(); err != nil {
		return err
	}
	return nil
}
