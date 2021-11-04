package hyperv

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

func resourceHyperVFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceHyperVFileCreate,
		Read:   resourceHyperVFileRead,
		Update: resourceHyperVFileUpdate,
		Delete: resourceHyperVFileDelete,

		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},

			"source": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"dirname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"exists": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"creationtime": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"lastwritetime": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		//CustomizeDiff: customizeDiffForFile,
	}
}

// func customizeDiffForFile(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
// 	path := diff.Get("path").(string)

// 	if _, err := os.Stat(path); err != nil {
// 		if os.IsNotExist(err) {
// 			// file does not exist
// 			diff.SetNewComputed("exists")
// 			return nil
// 		} else {
// 			// other error
// 			return err
// 		}
// 	}

// 	return nil
// }

func resourceHyperVFileCreate(d *schema.ResourceData, meta interface{}) (err error) {

	var path, source string
	log.Printf("[INFO][hyperv][create] creating hyperv File: %#v", d)
	c := meta.(*api.HypervClient)

	raw_path, ok := d.GetOk("path")
	if ok {
		path = raw_path.(string)
	} else {
		return fmt.Errorf("[ERROR][hyperv][create] path argument is required")
	}

	raw_source, ok := d.GetOk("source")
	if ok {
		source = raw_source.(string)
	} else {
		return fmt.Errorf("[ERROR][hyperv][create] source argument is required")
	}

	err = c.CreateOrUpdateFile(path, source)

	if err != nil {
		return err
	}

	d.SetId(path)

	log.Printf("[INFO][hyperv][create] created hyperv File: %#v", d)

	return resourceHyperVFileRead(d, meta)
}

func resourceHyperVFileRead(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[INFO][hyperv][read] reading hyperv File: %#v", d)
	c := meta.(*api.HypervClient)

	var path string
	raw_path, ok := d.GetOk("path")
	if ok {
		path = raw_path.(string)
	} else {
		return fmt.Errorf("[ERROR][hyperv][read] path argument is required")
	}

	file, err := c.GetFile(path)
	if err != nil {
		return err
	}

	d.SetId(path)
	//d.Set("path", file.Path)

	if file.Path != "" {
		log.Printf("[INFO][hyperv][read] unable to retrieved File: %+v", path)
		d.Set("exists", false)
	} else {
		log.Printf("[INFO][hyperv][read] retrieved File: %+v", path)
		d.Set("name", file.Name)
		d.Set("size", file.Size)
		d.Set("dirname", file.DirName)
		d.Set("exists", file.Exists)
		d.Set("creationtime", file.CreationTime)
		d.Set("lastwritetime", file.LastWriteTime)
	}

	log.Printf("[INFO][hyperv][read] read hyperv File: %#v", d)

	return nil
}

func resourceHyperVFileUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	var path, source string
	log.Printf("[INFO][hyperv][update] updating hyperv File: %#v", d)
	c := meta.(*api.HypervClient)

	raw_path, ok := d.GetOk("path")
	if ok {
		path = raw_path.(string)
	} else {
		return fmt.Errorf("[ERROR][hyperv][update] path argument is required")
	}

	raw_source, ok := d.GetOk("source")
	if ok {
		source = raw_source.(string)
	} else {
		return fmt.Errorf("[ERROR][hyperv][update] source argument is required")
	}

	exists := (d.Get("exists")).(bool)

	if !exists || d.HasChange("path") || d.HasChange("source") {
		//delete it as its changed
		err = c.CreateOrUpdateFile(path, source)

		if err != nil {
			return err
		}
	}

	log.Printf("[INFO][hyperv][update] updated hyperv File: %#v", d)

	return resourceHyperVFileRead(d, meta)
}

func resourceHyperVFileDelete(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[INFO][hyperv][delete] deleting hyperv File: %#v", d)

	c := meta.(*api.HypervClient)

	var path string
	raw_path, ok := d.GetOk("path")
	if ok {
		path = raw_path.(string)
	} else {
		return fmt.Errorf("[ERROR][hyperv][delete] path argument is required")
	}

	err = c.DeleteFile(path)

	if err != nil {
		return err
	}

	log.Printf("[INFO][hyperv][delete] deleted hyperv File: %#v", d)
	return nil
}
