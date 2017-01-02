// Package driver coordinates distributed execution.
package driver

import (
	"github.com/chrislusf/gleam/distributed/plan"
	"github.com/chrislusf/gleam/flow"
	"github.com/chrislusf/gleam/pb"
)

func (fcd *FlowContextDriver) GetTaskGroupStatus(taskGroup *plan.TaskGroup) *pb.FlowExecutionStatus_TaskGroup {
	for _, status := range fcd.status.TaskGroups {
		if len(taskGroup.Tasks) == len(status.TaskIds) {
			if int32(taskGroup.Tasks[0].Id) == status.TaskIds[0] {
				return status
			}
		}
	}
	return nil
}

func (fcd *FlowContextDriver) logExecutionPlan(fc *flow.FlowContext) {

	for _, step := range fc.Steps {
		var parentIds, taskIds, inputDatasetIds []int32
		for _, inputDataset := range step.InputDatasets {
			parentIds = append(parentIds, int32(inputDataset.Step.Id))
			inputDatasetIds = append(inputDatasetIds, int32(inputDataset.Id))
		}
		for _, task := range step.Tasks {
			taskIds = append(taskIds, int32(task.Id))
		}
		outputDatasetId := int32(0)
		if step.OutputDataset != nil {
			outputDatasetId = int32(step.OutputDataset.Id)
		}
		fcd.status.Steps = append(
			fcd.status.Steps,
			&pb.FlowExecutionStatus_Step{
				Id:              int32(step.Id),
				Name:            step.Name,
				ParentIds:       parentIds,
				TaskIds:         taskIds,
				InputDatasetId:  inputDatasetIds,
				OutputDatasetId: outputDatasetId,
			},
		)
	}

	for _, step := range fc.Steps {
		for _, task := range step.Tasks {
			fcd.status.Tasks = append(
				fcd.status.Tasks,
				&pb.FlowExecutionStatus_Task{
					StepId: int32(step.Id),
					Id:     int32(task.Id),
				},
			)
		}
	}

	for _, stepGroup := range fcd.stepGroups {
		var stepIds []int32
		for _, step := range stepGroup.Steps {
			stepIds = append(stepIds, int32(step.Id))
		}
		fcd.status.StepGroups = append(
			fcd.status.StepGroups,
			&pb.FlowExecutionStatus_StepGroup{
				StepIds: stepIds,
			},
		)
	}

	for _, dataset := range fc.Datasets {
		for _, shard := range dataset.Shards {
			fcd.status.DatasetShards = append(
				fcd.status.DatasetShards,
				&pb.FlowExecutionStatus_DatasetShard{
					DatasetId: int32(dataset.Id),
					Id:        int32(shard.Id),
				},
			)
		}
		var stepIds []int32
		for _, step := range dataset.ReadingSteps {
			stepIds = append(stepIds, int32(step.Id))
		}
		fcd.status.Datasets = append(
			fcd.status.Datasets,
			&pb.FlowExecutionStatus_Dataset{
				Id:             int32(dataset.Id),
				StepId:         int32(dataset.Step.Id),
				ReadingStepIds: stepIds,
			},
		)
	}

	for _, taskGroup := range fcd.taskGroups {
		var taskIds []int32
		for _, task := range taskGroup.Tasks {
			taskIds = append(taskIds, int32(task.Id))
		}
		statusTaskGroup := &pb.FlowExecutionStatus_TaskGroup{
			TaskIds: taskIds,
		}
		fcd.status.TaskGroups = append(
			fcd.status.TaskGroups,
			statusTaskGroup,
		)
	}

}