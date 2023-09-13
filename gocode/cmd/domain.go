/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/

package cmd

import (
	"github.com/spf13/cobra"
	"go.olapie.com/tools/gocode/domain"
	"go.olapie.com/utils"
)

// domainCmd represents the domain command
var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Generate domain code",
	Long: `Generate domain code with XML model. E.g. 
<?xml version="1.0" encoding="UTF-8" ?>
<model jsonNaming="SnakeCase" bsonNaming="CamelCase">
    <import>time</import>
    <import>context</import>
    <alias type="string">OrderTitle</alias>
    <simpletype type="int64">OrderID</simpletype>
    <simpletype type="int64">ItemID</simpletype>
    <simpletype type="string">Decimal</simpletype>
    <struct name="Item" json="true">
        <field type="ItemID">ID</field>
        <field type="Decimal">Price</field>
    </struct>
    <struct name="Order" json="true">
        <field type="OrderID">ID</field>
        <field type="time.Time">CreatedAt</field>
    </struct>
    <struct name="ItemEntity" json="true">
        <field type="ItemID">ID</field>
        <field type="Decimal">Price</field>
    </struct>
    <entity name="OrderEntity" json="true" bson="true">
        <field type="OrderID" bson="_id" readonly="true">ID</field>
        <field type="OrderTitle">Title</field>
        <field type="[]ItemID">ItemIDs</field>
        <field type="time.Time">ExpectedShipmentTime</field>
        <field type="time.Time">CreatedAt</field>
        <method>TotalPrice() Decimal</method>
    </entity>
    <interface name="ItemRepo">
        <method>Get(ctx context.Context, id ItemID)(*ItemEntity, error)</method>
    </interface>
    <interface name="OrderRepo">
        <method>Get(ctx context.Context, id OrderID)(*OrderEntity, error)</method>
    </interface>
    <interface name="OrderService">
        <method>Get(ctx context.Context, id OrderID)(*Order, error)</method>
    </interface>
</model>


`,
	Run: func(cmd *cobra.Command, args []string) {
		inputFilename := utils.MustGet(cmd.Flags().GetString("filename"))
		outputFilename := utils.MustGet(cmd.Flags().GetString("output"))
		domain.Generate(inputFilename, outputFilename)
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// domainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	domainCmd.Flags().StringP("filename", "f", "domain.xml", "Domain model XML filename")
	domainCmd.Flags().StringP("output", "o", "domain.gen.go", "Generated filename")
}
