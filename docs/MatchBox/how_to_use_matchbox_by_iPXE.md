# How to use Matchbox by iPXE
This tutorial show how to add label and reach to matchbox with KVM.

Refer [official doc](https://coreos.com/matchbox/docs/latest/network-setup.html)

**iPXE**
```
http://matchbox.example.com:port/boot.ipxe?lable=value
```

**Example**
```
virsh net-dumpxml upi

<network connections='1'>
  <name>upi</name>
  ...
  <ip family='ipv4' address='192.168.222.1' prefix='24'>
    <dhcp>
      <range start='192.168.222.2' end='192.168.222.254'/>
      ....
      <bootp file='http://matchbox.example.com:8080/boot.ipxe'/> <====
    </dhcp>
  </ip>
</network>

```

**Add bootp using terraform libvirt provider for KVM**

Create boot.xsl
```
<?xml version="1.0" ?>
<xsl:stylesheet version="1.0"
                xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
 <xsl:output omit-xml-declaration="yes" indent="yes"/>
 <!-- Identity transform -->
   <xsl:template match="@* | node()">
      <xsl:copy>
         <xsl:apply-templates select="@* | node()"/>
      </xsl:copy>
   </xsl:template>
  <xsl:template match="/network/ip/dhcp/range">
       <xsl:copy-of select="."/>
       <bootp file="http://matchbox.example.com:8080/boot.ipxe" />
  </xsl:template>

</xsl:stylesheet>
```

Create network.tf
```
resource "libvirt_network" "ocp_network" {
  ....
   xml {
       xslt = "${file("bootp.xsl")}"
   }
  ....
}
```