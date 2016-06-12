package null

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func repro() *schema.Resource {
	return &schema.Resource{
		Create: reproCreate,
		Read:   reproRead,
		Update: reproUpdate,
		Delete: reproDelete,

		Schema: map[string]*schema.Schema{
			"main_list": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"inner_string_set": &schema.Schema{
							Type:     schema.TypeSet,
							Required: true,
							Set:      schema.HashString,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"inner_bool": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"inner_int": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"main_set": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"inner_string_set": &schema.Schema{
							Type:     schema.TypeSet,
							Required: true,
							Set:      schema.HashString,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"inner_bool": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"inner_int": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"main_string_set": &schema.Schema{
				Type:     schema.TypeSet,
				Set:      schema.HashString,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"main_bool": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"main_int": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func reproCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(fmt.Sprintf("%d", rand.Int()))

	return reproRead(d, meta)
}

func reproRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func reproUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] UPDATE main_bool: %t", d.Get("main_bool").(bool))
	log.Printf("[DEBUG] UPDATE main_int: %d", d.Get("main_int").(int))
	l := d.Get("main_string_set").(*schema.Set).List()
	log.Printf("[DEBUG] UPDATE main_string_set (List): %q", l)

	// main list
	ml := d.Get("main_list").([]interface{})
	if len(ml) > 0 {
		m := ml[0].(map[string]interface{})

		log.Printf("[DEBUG] UPDATE main_list.inner_string_set: %q", m["inner_string_set"].(*schema.Set).List())
		log.Printf("[DEBUG] UPDATE main_list.inner_bool: %t", m["inner_bool"].(bool))
		log.Printf("[DEBUG] UPDATE main_list.inner_int: %d", m["inner_int"].(int))
	}

	// main set
	ms := d.Get("main_set").(*schema.Set).List()
	if len(ms) > 0 {
		m := ms[0].(map[string]interface{})

		log.Printf("[DEBUG] UPDATE main_set.inner_string_set: %q", m["inner_string_set"].(*schema.Set).List())
		log.Printf("[DEBUG] UPDATE main_set.inner_bool: %t", m["inner_bool"].(bool))
		log.Printf("[DEBUG] UPDATE main_set.inner_int: %d", m["inner_int"].(int))
	}

	return nil
}

func reproDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
