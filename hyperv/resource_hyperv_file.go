package hyperv

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

func resourceHyperVFile() *schema.Resource {
	return &schema.Resource{ // BUAK: Terraform hooks to perform CRUD operations on the resource
		Create: resourceHyperVFileCreate,
		Read:   resourceHyperVFileRead,
		Update: resourceHyperVFileUpdate,
		Delete: resourceHyperVFileDelete,

		Schema: map[string]*schema.Schema{ // BUAK: Schema are variables for the Terraform resource definition (in *.tf files)
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
	}
}

func resourceHyperVFileCreate(d *schema.ResourceData, meta interface{}) (err error) {
	// BUAK: Function creates a file on Hyper-V by downloading / copying from source
	var path, source string
	log.Printf("[INFO][hyperv][create] creating hyperv File: %#v", d)
	c := meta.(*api.HypervClient) // BUAK: Reference to api module (the thing with .ps1 scripts)

	raw_path, ok := d.GetOk("path") // BUAK: GetOK() has 2 return values (obj & bool)
	if ok {
		path = raw_path.(string) // BUAK: type cast to string
	} else {
		return fmt.Errorf("[ERROR][hyperv][create] path argument is required") // BUAK: If not ok return exception
	}

	raw_source, ok := d.GetOk("source") // BUAK: see path (above)
	if ok {
		source = raw_source.(string)
	} else {
		return fmt.Errorf("[ERROR][hyperv][create] source argument is required")
	}

	err = c.CreateOrUpdateFile(path, source) // BUAK: Call CreateOrUpdatefile() on api module (method def in api/file.go)

	if err != nil {
		return err
	}

	d.SetId(path) // BUAK: Set Terraform resource ID equal to path

	log.Printf("[INFO][hyperv][create] created hyperv File: %#v", d)

	return resourceHyperVFileRead(d, meta) // BUAK: Refresh state & exit
}

func resourceHyperVFileRead(d *schema.ResourceData, meta interface{}) (err error) {
	// BUAK: Function that collects infos about the file
	log.Printf("[INFO][hyperv][read] reading hyperv File: %#v", d)
	c := meta.(*api.HypervClient) // BUAK: Reference to api module (the thing with .ps1 scripts)

	var path string
	raw_path, ok := d.GetOk("path") // BUAK: GetOK() has 2 return values (obj & bool)
	if ok {
		path = raw_path.(string) // BUAK: type cast to string
	} else {
		return fmt.Errorf("[ERROR][hyperv][read] path argument is required") // BUAK: If not ok return exception
	}

	file, err := c.GetFile(path)
	if err != nil {
		return err
	}

	d.SetId(file.Path) // BUAK: Set Terraform resource ID equal to the path dedected by GetFile.
						// If the file does not exist, set ID to empty and tell terraform that the file was deleted on hyperv
	if file.Path == "" {
		log.Printf("[INFO][hyperv][read] unable to retrieved File: %+v", path)
		d.Set("exists", false)
	} else {
		// BUAK: Save info to Terraform state to be reused
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
	// BUAK: Function to update a file on Hyper-V
	// BUAK: Check resourceHyperVFileCreate() function
	var path, source string
	log.Printf("[INFO][hyperv][update] updating hyperv File: %#v", d)
	c := meta.(*api.HypervClient) // BUAK: Reference to api module (the thing with .ps1 scripts)

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

	exists := (d.Get("exists")).(bool) // BUAK: Check if the file still exists (it might have been manually deleted on hyperv)

	if !exists || d.HasChange("path") || d.HasChange("source") { // BUAK: IF it does not exist OR desired destination path changed OR source changed
		err = c.CreateOrUpdateFile(path, source)

		if err != nil {
			return err
		}
	}

	log.Printf("[INFO][hyperv][update] updated hyperv File: %#v", d)

	return resourceHyperVFileRead(d, meta)
}

func resourceHyperVFileDelete(d *schema.ResourceData, meta interface{}) (err error) {
	// BUAK: Function to delete file on Hyper-V
	// BUAK: Check resourceHyperVFileCreate() function
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
