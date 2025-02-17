// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cli

import (
	"fmt"

	"github.com/urfave/cli"
	enumspb "go.temporal.io/api/enums/v1"
)

func newAdminWorkflowCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "show",
			Aliases: []string{"show"},
			Usage:   "show workflow history from database",
			Flags: append(getDBFlags(),
				// v2 history events
				cli.StringFlag{
					Name:  FlagTreeID,
					Usage: "TreeId",
				},
				cli.StringFlag{
					Name:  FlagBranchID,
					Usage: "BranchId",
				},
				cli.StringFlag{
					Name:  FlagOutputFilenameWithAlias,
					Usage: "output file",
				},
				// support mysql query
				cli.IntFlag{
					Name:  FlagShardIDWithAlias,
					Usage: "ShardId",
				}),
			Action: func(c *cli.Context) {
				AdminShowWorkflow(c)
			},
		},
		{
			Name:    "describe",
			Aliases: []string{"desc"},
			Usage:   "Describe internal information of workflow execution",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  FlagWorkflowIDWithAlias,
					Usage: "WorkflowId",
				},
				cli.StringFlag{
					Name:  FlagRunIDWithAlias,
					Usage: "RunId",
				},
			},
			Action: func(c *cli.Context) {
				AdminDescribeWorkflow(c)
			},
		},
		{
			Name:    "refresh_tasks",
			Aliases: []string{"rt"},
			Usage:   "Refreshes all the tasks of a workflow",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  FlagWorkflowIDWithAlias,
					Usage: "WorkflowId",
				},
				cli.StringFlag{
					Name:  FlagRunIDWithAlias,
					Usage: "RunId",
				},
			},
			Action: func(c *cli.Context) {
				AdminRefreshWorkflowTasks(c)
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"del"},
			Usage:   "Delete current workflow execution and the mutableState record",
			Flags: append(
				getDBAndESFlags(),
				cli.StringFlag{
					Name:  FlagWorkflowIDWithAlias,
					Usage: "WorkflowId",
				},
				cli.StringFlag{
					Name:  FlagRunIDWithAlias,
					Usage: "RunId",
				},
				cli.BoolFlag{
					Name:  FlagSkipErrorModeWithAlias,
					Usage: "skip errors",
				}),
			Action: func(c *cli.Context) {
				AdminDeleteWorkflow(c)
			},
		},
	}
}

func newAdminShardManagementCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "describe",
			Aliases: []string{"d"},
			Usage:   "Describe shard by Id",
			Flags: append(
				getDBFlags(),
				cli.IntFlag{
					Name:  FlagShardID,
					Usage: "The Id of the shard to describe",
				},
				cli.StringFlag{
					Name:  FlagTargetCluster,
					Value: "active",
					Usage: "Temporal cluster to use",
				},
			),
			Action: func(c *cli.Context) {
				AdminDescribeShard(c)
			},
		},
		{
			Name:    "describe_task",
			Aliases: []string{"dt"},
			Usage:   "Describe a task based on task Id, task type, shard Id and task visibility timestamp",
			Flags: append(
				getDBFlags(),
				cli.IntFlag{
					Name:  FlagShardID,
					Usage: "The ID of the shard",
				},
				cli.IntFlag{
					Name:  FlagTaskID,
					Usage: "The ID of the timer task to describe",
				},
				cli.StringFlag{
					Name:  FlagTaskType,
					Value: "transfer",
					Usage: "Task type: transfer (default), timer, replication",
				},
				cli.Int64Flag{
					Name:  FlagTaskVisibilityTimestamp,
					Usage: "Task visibility timestamp in nano",
				},
				cli.StringFlag{
					Name:  FlagTargetCluster,
					Value: "active",
					Usage: "Temporal cluster to use",
				},
			),
			Action: func(c *cli.Context) {
				AdminDescribeTask(c)
			},
		},
		{
			Name:  "list_tasks",
			Usage: "List tasks for given shard Id and task type",
			Flags: append(append(
				getDBFlags(),
				flagsForPagination...),
				cli.StringFlag{
					Name:  FlagTargetCluster,
					Value: "active",
					Usage: "Temporal cluster to use",
				},
				cli.IntFlag{
					Name:  FlagShardID,
					Usage: "The ID of the shard",
				},
				cli.StringFlag{
					Name:  FlagTaskType,
					Value: "transfer",
					Usage: "Task type: transfer (default), timer, replication",
				},
				cli.StringFlag{
					Name:  FlagMinVisibilityTimestamp,
					Value: "2020-01-01T00:00:00+00:00",
					Usage: "Task visibility min timestamp. Supported formats are '2006-01-02T15:04:05+07:00', raw UnixNano and " +
						"time range (N<duration>), where 0 < N < 1000000 and duration (full-notation/short-notation) can be second/s, " +
						"minute/m, hour/h, day/d, week/w, month/M or year/y. For example, '15minute' or '15m' implies last 15 minutes.",
				},
				cli.StringFlag{
					Name:  FlagMaxVisibilityTimestamp,
					Value: "2035-01-01T00:00:00+00:00",
					Usage: "Task visibility max timestamp. Supported formats are '2006-01-02T15:04:05+07:00', raw UnixNano and " +
						"time range (N<duration>), where 0 < N < 1000000 and duration (full-notation/short-notation) can be second/s, " +
						"minute/m, hour/h, day/d, week/w, month/M or year/y. For example, '15minute' or '15m' implies last 15 minutes.",
				},
			),
			Action: func(c *cli.Context) {
				AdminListTasks(c)
			},
		},
		{
			Name:    "close_shard",
			Aliases: []string{"clsh"},
			Usage:   "close a shard given a shard id",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  FlagShardID,
					Usage: "ShardId for the temporal cluster to manage",
				},
			},
			Action: func(c *cli.Context) {
				AdminShardManagement(c)
			},
		},
		{
			Name:    "remove_task",
			Aliases: []string{"rmtk"},
			Usage:   "remove a task based on shardId, task type, taskId, and task visibility timestamp",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  FlagShardID,
					Usage: "shardId",
				},
				cli.Int64Flag{
					Name:  FlagTaskID,
					Usage: "taskId",
				},
				cli.StringFlag{
					Name:  FlagTaskType,
					Value: "transfer",
					Usage: "Task type: transfer (default), timer, replication",
				},
				cli.Int64Flag{
					Name:  FlagTaskVisibilityTimestamp,
					Usage: "task visibility timestamp in nano (required for removing timer task)",
				},
			},
			Action: func(c *cli.Context) {
				AdminRemoveTask(c)
			},
		},
	}
}

func newAdminMembershipCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "list_gossip",
			Usage: "List ringpop membership items",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  FlagClusterMembershipRole,
					Value: "all",
					Usage: "Membership role filter: all (default), frontend, history, matching, worker",
				},
			},
			Action: func(c *cli.Context) {
				AdminListGossipMembers(c)
			},
		},
		{
			Name:  "list_db",
			Usage: "List cluster membership items",
			Flags: append(
				getDBFlags(),
				cli.StringFlag{
					Name:  FlagHeartbeatedWithin,
					Value: "15m",
					Usage: "Filter by last heartbeat date time. Supported formats are '2006-01-02T15:04:05+07:00', raw UnixNano and " +
						"time range (N<duration>), where 0 < N < 1000000 and duration (full-notation/short-notation) can be second/s, " +
						"minute/m, hour/h, day/d, week/w, month/M or year/y. For example, '15minute' or '15m' implies last 15 minutes.",
				},
				cli.StringFlag{
					Name:  FlagClusterMembershipRole,
					Value: "all",
					Usage: "Membership role filter: all (default), frontend, history, matching, worker",
				},
			),
			Action: func(c *cli.Context) {
				AdminListClusterMembership(c)
			},
		},
	}
}

func newAdminHistoryHostCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "describe",
			Aliases: []string{"desc"},
			Usage:   "Describe internal information of history host",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  FlagWorkflowIDWithAlias,
					Usage: "WorkflowId",
				},
				cli.StringFlag{
					Name:  FlagHistoryAddressWithAlias,
					Usage: "History Host address(IP:PORT)",
				},
				cli.IntFlag{
					Name:  FlagShardIDWithAlias,
					Usage: "ShardId",
				},
				cli.BoolFlag{
					Name:  FlagPrintFullyDetailWithAlias,
					Usage: "Print fully detail",
				},
			},
			Action: func(c *cli.Context) {
				AdminDescribeHistoryHost(c)
			},
		},
		{
			Name:    "get_shardid",
			Aliases: []string{"gsh"},
			Usage:   "Get shardId for a namespaceId and workflowId combination",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  FlagNamespaceID,
					Usage: "NamespaceId",
				},
				cli.StringFlag{
					Name:  FlagWorkflowIDWithAlias,
					Usage: "WorkflowId",
				},
				cli.IntFlag{
					Name:  FlagNumberOfShards,
					Usage: "NumberOfShards for the temporal cluster(see config for numHistoryShards)",
				},
			},
			Action: func(c *cli.Context) {
				AdminGetShardID(c)
			},
		},
	}
}

func newAdminNamespaceCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "list",
			Usage: "List namespaces",
			Flags: append(getDBFlags(), getFlagsForList()...),
			Action: func(c *cli.Context) {
				AdminListNamespaces(c)
			},
		},
		{
			Name:    "register",
			Aliases: []string{"re"},
			Usage:   "Register workflow namespace",
			Flags:   adminRegisterNamespaceFlags,
			Action: func(c *cli.Context) {
				newNamespaceCLI(c, true).RegisterNamespace(c)
			},
		},
		{
			Name:    "update",
			Aliases: []string{"up", "u"},
			Usage:   "Update existing workflow namespace",
			Flags:   adminUpdateNamespaceFlags,
			Action: func(c *cli.Context) {
				newNamespaceCLI(c, true).UpdateNamespace(c)
			},
		},
		{
			Name:    "describe",
			Aliases: []string{"desc"},
			Usage:   "Describe existing workflow namespace",
			Flags:   adminDescribeNamespaceFlags,
			Action: func(c *cli.Context) {
				newNamespaceCLI(c, true).DescribeNamespace(c)
			},
		},
		{
			Name:    "get_namespaceidorname",
			Aliases: []string{"getdn"},
			Usage:   "Get namespaceId or namespace",
			Flags: append(getDBFlags(),
				cli.StringFlag{
					Name:  FlagNamespace,
					Usage: "Namespace",
				},
				cli.StringFlag{
					Name:  FlagNamespaceID,
					Usage: "Namespace Id(uuid)",
				}),
			Action: func(c *cli.Context) {
				AdminGetNamespaceIDOrName(c)
			},
		},
	}
}

func newAdminElasticSearchCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "catIndex",
			Aliases: []string{"cind"},
			Usage:   "Cat Indices on Elasticsearch",
			Flags:   getESFlags(false),
			Action: func(c *cli.Context) {
				AdminCatIndices(c)
			},
		},
		{
			Name:    "index",
			Aliases: []string{"ind"},
			Usage:   "Index docs on Elasticsearch",
			Flags: append(
				getESFlags(true),
				cli.StringFlag{
					Name:  FlagInputFileWithAlias,
					Usage: "Input file of indexerspb.Message in json format, separated by newline",
				},
				cli.IntFlag{
					Name:  FlagBatchSizeWithAlias,
					Usage: "Optional batch size of actions for bulk operations",
					Value: 10,
				},
			),
			Action: func(c *cli.Context) {
				AdminIndex(c)
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"del"},
			Usage:   "Delete docs on Elasticsearch",
			Flags: append(
				getESFlags(true),
				cli.StringFlag{
					Name: FlagInputFileWithAlias,
					Usage: "Input file name. Redirect temporal wf list result (with table format) to a file and use as delete input. " +
						"First line should be table header like WORKFLOW TYPE | WORKFLOW ID | RUN ID | ...",
				},
				cli.IntFlag{
					Name:  FlagBatchSizeWithAlias,
					Usage: "Optional batch size of actions for bulk operations",
					Value: 1000,
				},
				cli.IntFlag{
					Name:  FlagRPS,
					Usage: "Optional batch request rate per second",
					Value: 30,
				},
			),
			Action: func(c *cli.Context) {
				AdminDelete(c)
			},
		},
		{
			Name:    "report",
			Aliases: []string{"rep"},
			Usage:   "Generate Report by Aggregation functions on Elasticsearch",
			Flags: append(
				getESFlags(true),
				cli.StringFlag{
					Name:  FlagListQuery,
					Usage: "SQL query of the report",
				},
				cli.StringFlag{
					Name:  FlagOutputFormat,
					Usage: "Additional output format (html or csv)",
				},
				cli.StringFlag{
					Name:  FlagOutputFilename,
					Usage: "Additional output filename with path",
				},
			),
			Action: func(c *cli.Context) {
				GenerateReport(c)
			},
		},
	}
}

func newAdminTaskQueueCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "describe",
			Aliases: []string{"desc"},
			Usage:   "Describe pollers and status information of task queue",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  FlagTaskQueueWithAlias,
					Usage: "TaskQueue description",
				},
				cli.StringFlag{
					Name:  FlagTaskQueueTypeWithAlias,
					Value: "workflow",
					Usage: "Optional TaskQueue type [workflow|activity]",
				},
			},
			Action: func(c *cli.Context) {
				AdminDescribeTaskQueue(c)
			},
		},
		{
			Name:  "list_tasks",
			Usage: "List tasks of a task queue",
			Flags: append(append(append(getDBFlags(), flagsForExecution...),
				flagsForPagination...),
				cli.StringFlag{
					Name:  FlagNamespaceID,
					Usage: "Namespace Id",
				},
				cli.StringFlag{
					Name:  FlagTaskQueueType,
					Value: "activity",
					Usage: "Taskqueue type: activity, workflow",
				},
				cli.StringFlag{
					Name:  FlagTaskQueue,
					Usage: "Taskqueue name",
				},
				cli.Int64Flag{
					Name:  FlagMinReadLevel,
					Usage: "Lower bound of read level",
				},
				cli.Int64Flag{
					Name:  FlagMaxReadLevel,
					Usage: "Upper bound of read level",
				},
			),
			Action: func(c *cli.Context) {
				AdminListTaskQueueTasks(c)
			},
		},
	}
}

func newAdminClusterCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "add-search-attributes",
			Aliases: []string{"asa"},
			Usage:   "Add custom search attributes",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:     FlagSkipSchemaUpdate,
					Usage:    "Skip Elasticsearch index schema update (only register in metadata)",
					Required: false,
				},
				cli.StringFlag{
					Name:   FlagIndex,
					Usage:  "Elasticsearch index name (optional)",
					Hidden: true, // don't show it for now
				},
				cli.StringSliceFlag{
					Name:  FlagNameWithAlias,
					Usage: "Search attribute name (multiply values are supported)",
				},
				cli.StringSliceFlag{
					Name:  FlagTypeWithAlias,
					Usage: fmt.Sprintf("Search attribute type: %v (multiply values are supported)", allowedEnumValues(enumspb.IndexedValueType_name)),
				},
			},
			Action: func(c *cli.Context) {
				AdminAddSearchAttributes(c)
			},
		},
		{
			Name:    "remove-search-attributes",
			Aliases: []string{"rsa"},
			Usage:   "Remove custom search attributes metadata only (Elasticsearch index schema is not modified)",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   FlagIndex,
					Usage:  "Elasticsearch index name (optional)",
					Hidden: true, // don't show it for now
				},
				cli.StringSliceFlag{
					Name:  FlagNameWithAlias,
					Usage: "Search attribute name",
				},
			},
			Action: func(c *cli.Context) {
				AdminRemoveSearchAttributes(c)
			},
		},
		{
			Name:    "get-search-attributes",
			Aliases: []string{"gsa"},
			Usage:   "Show existing search attributes",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  FlagPrintJSONWithAlias,
					Usage: "Output in JSON format",
				},
				cli.StringFlag{
					Name:   FlagIndex,
					Usage:  "Elasticsearch index name (optional)",
					Hidden: true, // don't show it for now
				},
			},
			Action: func(c *cli.Context) {
				AdminGetSearchAttributes(c)
			},
		},
		{
			Name:    "describe",
			Aliases: []string{"d"},
			Usage:   "Describe cluster information",
			Action: func(c *cli.Context) {
				AdminDescribeCluster(c)
			},
		},
		{
			Name:    "metadata",
			Aliases: []string{"m"},
			Usage:   "Show cluster metadata",
			Action: func(c *cli.Context) {
				AdminClusterMetadata(c)
			},
		},
	}
}

func newAdminDLQCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "read",
			Aliases: []string{"r"},
			Usage:   "Read DLQ Messages",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  FlagDLQTypeWithAlias,
					Usage: "Type of DLQ to manage. (Options: namespace, history)",
				},
				cli.StringFlag{
					Name:  FlagCluster,
					Usage: "Source cluster",
				},
				cli.IntFlag{
					Name:  FlagShardIDWithAlias,
					Usage: "ShardId",
				},
				cli.IntFlag{
					Name:  FlagMaxMessageCountWithAlias,
					Usage: "Max message size to fetch",
				},
				cli.IntFlag{
					Name:  FlagLastMessageID,
					Usage: "The upper boundary of the read message",
				},
				cli.StringFlag{
					Name:  FlagOutputFilenameWithAlias,
					Usage: "Output file to write to, if not provided output is written to stdout",
				},
			},
			Action: func(c *cli.Context) {
				AdminGetDLQMessages(c)
			},
		},
		{
			Name:    "purge",
			Aliases: []string{"p"},
			Usage:   "Delete DLQ messages with equal or smaller ids than the provided task id",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  FlagDLQTypeWithAlias,
					Usage: "Type of DLQ to manage. (Options: namespace, history)",
				},
				cli.StringFlag{
					Name:  FlagCluster,
					Usage: "Source cluster",
				},
				cli.IntFlag{
					Name:  FlagShardIDWithAlias,
					Usage: "ShardId",
				},
				cli.IntFlag{
					Name:  FlagLastMessageID,
					Usage: "The upper boundary of the read message",
				},
			},
			Action: func(c *cli.Context) {
				AdminPurgeDLQMessages(c)
			},
		},
		{
			Name:    "merge",
			Aliases: []string{"m"},
			Usage:   "Merge DLQ messages with equal or smaller ids than the provided task id",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  FlagDLQTypeWithAlias,
					Usage: "Type of DLQ to manage. (Options: namespace, history)",
				},
				cli.StringFlag{
					Name:  FlagCluster,
					Usage: "Source cluster",
				},
				cli.IntFlag{
					Name:  FlagShardIDWithAlias,
					Usage: "ShardId",
				},
				cli.IntFlag{
					Name:  FlagLastMessageID,
					Usage: "The upper boundary of the read message",
				},
			},
			Action: func(c *cli.Context) {
				AdminMergeDLQMessages(c)
			},
		},
	}
}

func newDBCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "scan",
			Aliases: []string{"scan"},
			Usage:   "scan concrete executions in database and detect corruptions",
			Flags: append(getDBFlags(),
				cli.IntFlag{
					Name:  FlagLowerShardBound,
					Usage: "lower bound of shard to scan (inclusive)",
					Value: 0,
				},
				cli.IntFlag{
					Name:  FlagUpperShardBound,
					Usage: "upper bound of shard to scan (exclusive)",
					Value: 16384,
				},
				cli.IntFlag{
					Name:  FlagStartingRPS,
					Usage: "starting rps of database queries, rps will be increased to target over scale up seconds",
					Value: 100,
				},
				cli.IntFlag{
					Name:  FlagRPS,
					Usage: "target rps of database queries, target will be reached over scale up seconds",
					Value: 7000,
				},
				cli.IntFlag{
					Name:  FlagPageSize,
					Usage: "page size used to query db executions table",
					Value: 500,
				},
				cli.IntFlag{
					Name:  FlagConcurrency,
					Usage: "number of threads to handle scan",
					Value: 1000,
				},
				cli.IntFlag{
					Name:  FlagReportRate,
					Usage: "the number of shards which get handled between each emitting of progress",
					Value: 10,
				}),
			Action: func(c *cli.Context) {
				AdminDBScan(c)
			},
		},
		{
			Name:    "clean",
			Aliases: []string{"clean"},
			Usage:   "clean up corrupted workflows",
			Flags: append(getDBFlags(),
				cli.StringFlag{
					Name:  FlagInputDirectory,
					Usage: "the directory which contains corrupted workflow execution files from scan",
				},
				cli.IntFlag{
					Name:  FlagLowerShardBound,
					Usage: "lower bound of corrupt shard to handle (inclusive)",
					Value: 0,
				},
				cli.IntFlag{
					Name:  FlagUpperShardBound,
					Usage: "upper bound of shard to handle (exclusive)",
					Value: 16384,
				},
				cli.IntFlag{
					Name:  FlagStartingRPS,
					Usage: "starting rps of database queries, rps will be increased to target over scale up seconds",
					Value: 100,
				},
				cli.IntFlag{
					Name:  FlagRPS,
					Usage: "target rps of database queries, target will be reached over scale up seconds",
					Value: 7000,
				},
				cli.IntFlag{
					Name:  FlagConcurrency,
					Usage: "number of threads to handle clean",
					Value: 1000,
				},
				cli.IntFlag{
					Name:  FlagReportRate,
					Usage: "the number of shards which get handled between each emitting of progress",
					Value: 10,
				}),
			Action: func(c *cli.Context) {
				AdminDBClean(c)
			},
		},
	}
}

func newDecodeCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "proto",
			Usage: "Decode proto payload",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  FlagProtoType,
					Usage: "full name of proto type to decode to (i.e. temporal.server.api.persistence.v1.WorkflowExecutionInfo).",
				},
				cli.StringFlag{
					Name:  FlagHexData,
					Usage: "data in hex format (i.e. 0x0a243462613036633466...).",
				},
				cli.StringFlag{
					Name:  FlagHexFile,
					Usage: "file with data in hex format (i.e. 0x0a243462613036633466...).",
				},
				cli.StringFlag{
					Name:  FlagBinaryFile,
					Usage: "file with data in binary format.",
				},
			},
			Action: func(c *cli.Context) {
				AdminDecodeProto(c)
			},
		},
		{
			Name:  "base64",
			Usage: "Decode base64 payload",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  FlagBase64Data,
					Usage: "data in base64 format (i.e. anNvbi9wbGFpbg==).",
				},
				cli.StringFlag{
					Name:  FlagBase64File,
					Usage: "file with data in base64 format (i.e. anNvbi9wbGFpbg==).",
				},
			},
			Action: func(c *cli.Context) {
				AdminDecodeBase64(c)
			},
		},
	}
}
