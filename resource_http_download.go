package main

import (
        "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
        "io"
        "net/http"
        "os"
        "crypto/md5"
        "crypto/sha256"
        "encoding/hex"
        "fmt"
)

func httpDownload() *schema.Resource {
        return &schema.Resource{
                Create: resourceHttpDownloadCreate,
                Read:   resourceHttpDownloadRead,
                Update: resourceHttpDownloadUpdate,
                Delete: resourceHttpDownloadDelete,

                Schema: map[string]*schema.Schema{
                        "remote_url": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "filename": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "checksum_type": &schema.Schema{
                                Type:     schema.TypeString,
                                Optional: true,
                        },
                        "checksum": &schema.Schema{
                                Type:     schema.TypeString,
                                Optional: true,
                        },
                },
        }
}


func resourceHttpDownloadCreate(d *schema.ResourceData, m interface{}) error {
        remote_url := d.Get("remote_url").(string)
        filename := d.Get("filename").(string)

        // Download the file
        resp, err := http.Get(remote_url)
        if err != nil {
                return err
        }
        defer resp.Body.Close()

        // Create the file
        file, err := os.Create(filename)
        if err != nil {
                return fmt.Errorf("Could write to file: %s", err)
        }
        defer file.Close()

        // Write the body to file
        _, err = io.Copy(file, resp.Body)
        if err != nil {
                return fmt.Errorf("Something went wrong while writing the file: %s", err)
        }

        // Tell Terraform that everyting went fine
        d.SetId(filename)

        // Important: Tell Terraform to perform the read-validation
        return resourceHttpDownloadRead(d, m)
}

func resourceHttpDownloadRead(d *schema.ResourceData, m interface{}) error {
        filename := d.Get("filename").(string)
        checksum := d.Get("checksum").(string)
        checksumType := d.Get("checksum_type").(string)

        checksum_on_disk := ""
        var err error = nil

        switch checksumType {
        case "md5":
                checksum_on_disk, err = getMd5Hash(filename)
        case "sha256":
                checksum_on_disk, err = getSha256Hash(filename)
        }

        if err != nil {
            // If there was an error calculating the checksum from disk (e.g. the file was removed outside terraform)
            d.SetId("")
            return fmt.Errorf("Error while calculating checksum: %s", err)
        }
        if checksum_on_disk != checksum {
            // If the calculated checksum doesn't match with the provided checksum: reset the Terraform id & error out.
            d.SetId("")
            return fmt.Errorf("Checksum of the downloaded file was: %s and did not match the specified checksum: %s", checksum_on_disk, checksum)
        }
        // Everything matches: set the Terraform id & checksum
        d.SetId(filename)
        d.Set("checksum", checksum_on_disk)

        return nil
}



func resourceHttpDownloadUpdate(d *schema.ResourceData, m interface{}) error {
        return resourceHttpDownloadCreate(d, m)
        return nil
}

func resourceHttpDownloadDelete(d *schema.ResourceData, m interface{}) error {
        filename := d.Get("filename").(string)

        err := os.Remove(filename)
        if err != nil {
            return err
        }
       return nil
}

func getMd5Hash(filename string)(string, error) {
       var result []byte
       file, err := os.Open(filename)
       if err != nil {
         return hex.EncodeToString(result), err
       }
       defer file.Close()

       hash := md5.New()
       if _, err := io.Copy(hash, file); err != nil {
         return hex.EncodeToString(result), err
       }

       return hex.EncodeToString(hash.Sum(result)), nil
}

func getSha256Hash(filename string)(string, error) {
       var result []byte
       file, err := os.Open(filename)
       if err != nil {
         return hex.EncodeToString(result), err
       }
       defer file.Close()

       hash := sha256.New()
       if _, err := io.Copy(hash, file); err != nil {
         return hex.EncodeToString(result), err
       }

       return hex.EncodeToString(hash.Sum(result)), nil
}
