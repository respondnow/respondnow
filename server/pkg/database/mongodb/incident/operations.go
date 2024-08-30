package incident

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/respondnow/respondnow/server/config"
	"github.com/respondnow/respondnow/server/pkg/constant"
	"github.com/respondnow/respondnow/server/pkg/database/mongodb"
)

func (i *incidentOperator) Create(ctx context.Context, in Incident,
	opts ...*options.InsertOneOptions) (Incident, error) {
	// set or reset some fields
	now := time.Now().Unix()
	in.ID = primitive.NilObjectID
	in.CreatedAt = now
	in.UpdatedAt = &now
	// validation
	if err := i.Validate(&in); err != nil {
		return Incident{}, fmt.Errorf("invalid input, error : %s", err)
	}
	// save
	res, err := i.operator.Create(ctx, mongodb.IncidentCollection, in, opts...)
	if err != nil {
		return Incident{}, err
	}
	return i.GetByID(ctx, res.InsertedID)
}

func (i *incidentOperator) GetByID(ctx context.Context, id interface{},
	opts ...*options.FindOneOptions) (Incident, error) {
	out := Incident{}
	res, err := i.operator.Get(ctx, mongodb.IncidentCollection,
		bson.D{{Key: constant.ID, Value: id}}, opts...)
	if err != nil {
		return out, err
	}
	if err := res.Decode(&out); err != nil {
		return Incident{}, err
	}
	return out, nil
}

func (i *incidentOperator) Get(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) (Incident, error) {
	out := Incident{}
	res, err := i.operator.Get(ctx, mongodb.IncidentCollection, filter, opts...)
	if err != nil {
		return out, err
	}
	if err := res.Decode(&out); err != nil {
		return Incident{}, err
	}
	return out, nil
}

func (i *incidentOperator) List(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) ([]Incident, error) {
	out := make([]Incident, 0)
	cursor, err := i.operator.List(ctx, mongodb.IncidentCollection, filter, opts...)
	if err != nil {
		return out, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		item := Incident{}
		if err = cursor.Decode(&item); err != nil {
			return out, err
		}
		out = append(out, item)
	}
	return out, nil
}

func (i *incidentOperator) CountDocuments(ctx context.Context, filter interface{},
	opts ...*options.CountOptions) (int64, error) {
	return i.operator.CountDocuments(ctx, mongodb.IncidentCollection, filter, opts...)
}

func (i *incidentOperator) UpdateByID(ctx context.Context, in Incident,
	opts ...*options.UpdateOptions) (Incident, error) {
	// set or reset some fields
	now := time.Now().Unix()
	in.UpdatedAt = &now
	if in.Removed {
		in.RemovedAt = &now
	}
	// validation
	if err := i.Validate(&in); err != nil {
		return Incident{}, fmt.Errorf("invalid input, error : %s", err)
	}
	// save
	update := bson.D{}
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "name", Value: in.Name}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "description", Value: in.Description}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "tags", Value: in.Tags}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "severity",
		Value: in.Severity}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "status",
		Value: in.Status}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "active", Value: in.Active}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "summary", Value: in.Summary}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "comment", Value: in.Comment}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "services", Value: in.Services}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "environments", Value: in.Environments}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "functionalities",
		Value: in.Functionalities}}})

	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "roles", Value: in.Roles}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "stages", Value: in.Stages}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "timelines", Value: in.Timelines}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "channels", Value: in.Channels}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "conferenceDetails",
		Value: in.ConferenceDetails}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "attachments", Value: in.Attachments}}})

	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "updatedAt", Value: in.UpdatedAt}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "updatedBy", Value: in.UpdatedBy}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "removed", Value: in.Removed}}})
	update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "removedAt", Value: in.RemovedAt}}})
	_, err := i.operator.UpdateByID(ctx, mongodb.IncidentCollection, in.ID, update, opts...)
	if err != nil {
		return Incident{}, err
	}
	return i.GetByID(ctx, in.ID)
}

func (i *incidentOperator) CustomList(ctx context.Context, filter interface{},
	limit, skip int64) ([]Incident, error) {
	out := make([]Incident, 0)
	isLimitZero := false
	if limit == 0 {
		isLimitZero = true
	}
	for ok := true; ok; {
		ok = false
		max := limit
		if limit != constant.MaxDocumentInOneCall {
			if limit > constant.MaxDocumentInOneCall {
				max = constant.MaxDocumentInOneCall
			}
		}
		limit = limit - max
		if isLimitZero {
			max = constant.MaxDocumentInOneCall
		}
		if max == 0 {
			continue
		}
		cursor, err := i.operator.List(ctx, mongodb.IncidentCollection, filter,
			options.Find().SetLimit(max).SetSkip(skip).SetSort(bson.D{bson.E{Key: "createdAt", Value: -1}}))
		if err != nil {
			return out, err
		}
		defer cursor.Close(ctx)
		if isLimitZero {
			skip = skip + constant.MaxDocumentInOneCall
		} else {
			skip = skip + max
		}
		for cursor.Next(ctx) {
			ok = true
			item := Incident{}
			if err = cursor.Decode(&item); err != nil {
				return out, err
			}
			out = append(out, item)
		}
	}
	return out, nil
}

func (i *incidentOperator) BulkProcessWithSessionContext(ctx mongo.SessionContext, createList,
	updateList []Incident, opts ...*options.BulkWriteOptions) error {
	if len(createList) == 0 && len(updateList) == 0 {
		return nil
	}
	// set or reset some fields
	wModels := make([]mongo.WriteModel, 0)
	now := time.Now().Unix()
	for _, item := range createList {
		item.ID = primitive.NilObjectID
		item.CreatedAt = now
		// validation
		if err := i.Validate(&item); err != nil {
			return fmt.Errorf("invalid input, error : %s", err)
		}
		wModels = append(wModels, mongo.NewInsertOneModel().SetDocument(item))
	}
	for _, item := range updateList {
		now := now
		// set or reset some fields
		item.UpdatedAt = &now
		if item.Removed {
			item.RemovedAt = &now
		}
		// validation
		if err := i.Validate(&item); err != nil {
			return fmt.Errorf("invalid input, error : %s", err)
		}
		// save
		update := bson.D{}
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "name", Value: item.Name}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "description", Value: item.Description}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "tags", Value: item.Tags}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "severity",
			Value: item.Severity}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "status",
			Value: item.Status}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "active", Value: item.Active}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "summary", Value: item.Summary}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "services", Value: item.Services}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "environments",
			Value: item.Environments}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "functionalities",
			Value: item.Functionalities}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "roles", Value: item.Roles}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "stages", Value: item.Stages}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "channels", Value: item.Channels}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "conferenceDetails",
			Value: item.ConferenceDetails}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "attachments", Value: item.Attachments}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "updatedAt", Value: item.Services}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "updatedBy", Value: item.UpdatedBy}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "removed", Value: item.Removed}}})
		update = append(update, primitive.E{Key: "$set", Value: bson.D{{Key: "removedAt", Value: item.RemovedAt}}})
		wModels = append(wModels, mongo.NewUpdateOneModel().
			SetFilter(bson.D{primitive.E{Key: constant.ID, Value: item.ID}}).SetUpdate(update))
	}
	if _, err := i.operator.BulkWrite(ctx, mongodb.IncidentCollection, wModels, opts...); err != nil {
		return err
	}
	return nil
}

func (i *incidentOperator) WithDefaults(in *Incident) {}

func (i *incidentOperator) Validate(in *Incident) error {
	if in.Identifier == "" {
		return fmt.Errorf("missing identifier")
	}
	if in.Name == "" {
		return fmt.Errorf("missing name")
	}
	if in.AccountIdentifier == "" {
		return fmt.Errorf("missing account ID")
	}
	if in.ProjectIdentifier != "" && in.OrgIdentifier == "" {
		return fmt.Errorf("missing organization ID")
	}
	if in.Type == "" {
		return fmt.Errorf("missing incident type")
	}
	if in.Status == "" {
		return fmt.Errorf("missing incident status")
	}
	if in.Severity == "" {
		return fmt.Errorf("missing severity")
	}
	if in.Summary == "" && in.Description == "" {
		return fmt.Errorf("either summary or description must not be empty")
	}
	if in.IncidentChannel.Type == "" {
		return fmt.Errorf("missing incident channel type")
	}

	return nil
}

func (i *incidentOperator) GetIncidentTypes() []Type {
	resp := make([]Type, 0)
	if len(config.ServerConfig.IncidentTypes) > 0 {
		for _, incidentType := range config.ServerConfig.IncidentTypes {
			resp = append(resp, Type(incidentType))
		}
	} else {
		defaultSupportedTypes := []Type{
			Availability,
			Latency,
			Security,
			Other,
		}
		resp = append(resp, defaultSupportedTypes...)
	}

	return resp
}

func (i *incidentOperator) GetIncidentSeverities() []Severity {
	resp := make([]Severity, 0)
	if len(config.ServerConfig.Severities) > 0 {
		for severity := range config.ServerConfig.Severities {
			resp = append(resp, Severity(severity))
		}
	} else {
		defaultSupportedSeverities := []Severity{
			Severity0,
			Severity1,
			Severity2,
		}
		resp = append(resp, defaultSupportedSeverities...)
	}

	return resp
}

func (i *incidentOperator) GetIncidentAttachmentType() []AttachmentType {
	return []AttachmentType{
		Link,
	}
}

func (i *incidentOperator) GetIncidentStageStatuses() []Status {
	resp := make([]Status, 0)
	if len(config.ServerConfig.Statuses) > 0 {
		for _, status := range config.ServerConfig.Statuses {
			resp = append(resp, Status(status))
		}
	} else {
		defaultSupportedStatuses := []Status{
			Started,
			Acknowledged,
			Investigating,
			Identified,
			Mitigated,
			Resolved,
		}
		resp = append(resp, defaultSupportedStatuses...)
	}

	return resp
}

func (i *incidentOperator) GetIncidentRoles() []RoleType {
	resp := make([]RoleType, 0)
	if len(config.ServerConfig.Roles) > 0 {
		for role := range config.ServerConfig.Roles {
			resp = append(resp, RoleType(role))
		}
	} else {
		defaultSupportedRoles := []RoleType{
			IncidentCommander,
			CommunicationsLead,
		}
		resp = append(resp, defaultSupportedRoles...)
	}

	return resp
}
