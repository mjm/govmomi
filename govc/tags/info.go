/*
Copyright (c) 2018 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tags

import (
	"context"
	"flag"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/vmware/govmomi/govc/cli"
	"github.com/vmware/govmomi/govc/flags"
	"github.com/vmware/govmomi/vapi/tags"
)

type info struct {
	*flags.ClientFlag
	*flags.OutputFlag
	id bool
}

func init() {
	cli.Register("tags.info", &info{})
}

func (cmd *info) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.ClientFlag, ctx = flags.NewClientFlag(ctx)
	cmd.OutputFlag, ctx = flags.NewOutputFlag(ctx)
	cmd.ClientFlag.Register(ctx, f)
	cmd.OutputFlag.Register(ctx, f)
	f.BoolVar(&cmd.id, "i", false, "ID of category")
}

func (cmd *info) Process(ctx context.Context) error {
	if err := cmd.ClientFlag.Process(ctx); err != nil {
		return err
	}
	return nil
}
func (cmd *info) Usage() string {
	return "TAGNAME CATEGORYID or TAGID"
}

func (cmd *info) Description() string {
	return `Get tags info by tags' ID, or Get tags info for category by tag name and category ID. 

Examples:
  govc tags.info -i TAGID
  govc tags.info TAGNAME CATEGORYID `
}

type getTagInfo []tags.Tag

func (t getTagInfo) Write(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 2, 0, 2, ' ', 0)

	for _, item := range t {
		fmt.Fprintf(tw, "Name:\t%s\n", item.Name)
		fmt.Fprintf(tw, "  ID:\t%s\n", item.ID)
		fmt.Fprintf(tw, "  Description:\t%s\n", item.Description)
		fmt.Fprintf(tw, "  CategoryID:\t%s\n", item.CategoryID)
		fmt.Fprintf(tw, "  UsedBy: %s\n", item.UsedBy)
	}

	return tw.Flush()
}

func (cmd *info) Run(ctx context.Context, f *flag.FlagSet) error {

	return withClient(ctx, cmd.ClientFlag, func(c *tags.RestClient) error {
		var result getTagInfo
		if cmd.id {
			if f.NArg() != 1 {
				return flag.ErrHelp
			}
			id := f.Arg(0)
			tag, err := c.GetTag(ctx, id)
			if err != nil {
				return err
			}
			result = append(result, *tag)
		} else {
			if f.NArg() != 2 {
				return flag.ErrHelp
			}
			name := f.Arg(0)
			id := f.Arg(1)
			var err error
			result, err = c.GetTagByNameForCategory(ctx, name, id)
			if err != nil {
				return err
			}
		}
		return cmd.WriteResult(result)
	})
}
