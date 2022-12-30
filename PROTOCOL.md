# SMC 2x16 LCD

    $ lsusb -v -d 15d9:1133

    Bus 001 Device 010: ID 15d9:1133 Trust International B.V.
    Device Descriptor:
    bLength                18
    bDescriptorType         1
    bcdUSB               1.10
    bDeviceClass            0
    bDeviceSubClass         0
    bDeviceProtocol         0
    bMaxPacketSize0        32
    idVendor           0x15d9 Trust International B.V.
    idProduct          0x1133
    bcdDevice            1.00
    iManufacturer           1 SMC
    iProduct                2 SuperMicro LCD Display
    iSerial                 3 1.0.0
    bNumConfigurations      1
    Configuration Descriptor:
        bLength                 9
        bDescriptorType         2
        wTotalLength           41
        bNumInterfaces          1
        bConfigurationValue     1
        iConfiguration          0
        bmAttributes         0x80
        (Bus Powered)
        MaxPower              100mA
        Interface Descriptor:
        bLength                 9
        bDescriptorType         4
        bInterfaceNumber        0
        bAlternateSetting       0
        bNumEndpoints           2
        bInterfaceClass         3 Human Interface Device
        bInterfaceSubClass      0
        bInterfaceProtocol      0
        iInterface              0
            HID Device Descriptor:
            bLength                 9
            bDescriptorType        33
            bcdHID               1.11
            bCountryCode            0 Not supported
            bNumDescriptors         1
            bDescriptorType        34 Report
            wDescriptorLength      42
            Report Descriptors:
            ** UNAVAILABLE **
        Endpoint Descriptor:
            bLength                 7
            bDescriptorType         5
            bEndpointAddress     0x81  EP 1 IN
            bmAttributes            3
            Transfer Type              Interrupt
            Synch Type                 None
            Usage Type                 Data
            wMaxPacketSize     0x0010  1x 16 bytes
            bInterval               1
        Endpoint Descriptor:
            bLength                 7
            bDescriptorType         5
            bEndpointAddress     0x02  EP 2 OUT
            bmAttributes            3
            Transfer Type              Interrupt
            Synch Type                 None
            Usage Type                 Data
            wMaxPacketSize     0x0010  1x 16 bytes
            bInterval               1
    Device Status:     0x0001
    Self Powered

## Method

    download SMCLCD_v1.4.1.zip for Linux from https://www.supermicro.com/products/nfo/LCD.cfm

    decompile LCDMainUI.jar; found com.supermicro.lcd.SMCLCD class with
    references to native methods in libLIBSMC64.so.

    $ file libLIBSMC64.so
    libLIBSMC64.so: ELF 64-bit LSB shared object, x86-64, version 1 (SYSV),
    dynamically linked, with debug_info, not stripped

library not stripped; symbols dumped via nm to verify if other functions
possible but not used.

    int sl_init(void); /* returns handle */
    int sl_close(int handle); /* closes handle */
    int sl_clear(int handle);
    int sl_home(int handle);
    int sl_set_cursor(int handle, int type);
    int sl_get_key(int handle, int timeout);
    int sl_set_backlight(int handle, int state);
    int sl_move_cursor(int handle, int x, int y);
    int sl_putch(int handle, int c);
    int sl_printf(int handle, const char *s);
    int sl_mvprintf(int handle, int x, int y, const char *s);
    int sl_getch(int handle, int off, char *s);
    int sl_getline(int handle, int line, char *s);

## HID Reports

First byte is report ID; AA for input BB for output
Last byte appears to be checksum; 2s complement
Both reports are 16 bytes long
chksum choice seems odd; interrupt endpoints are guaranteed delivery unsure of
why this was added.
report followed by one command byte, and an optional subcommand

1.  LCD (`02`)

    1.  Control (`00`)

             00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15
            +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
            |BB|02|00|01|00|00|00|00|00|00|00|00|00|00|00|42|
            +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
            |ID|     |XX|            RESERVED            |CK|
            +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

             ID HID Report ID (BB)
             XX Command:
                01 Clear
                02 Home
                0C Cursor Off
                0D Block Cursor*
                0E Underline Cursor
                0F Block & Underline Cursor
                8X Move Cursor to Line 0 (+ Column)
                CX Move Cursor to Line 1 (+ Column)
             CK Checksum

             * Might be unintentional; decompiled java seems to indicate this
               is not a valid option.

    2.  Write (`02`)

             00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15
            +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
            |BB|02|02|62|61|6C|6C|73|00|00|00|00|00|00|00|33|
            +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
            |ID|     |               DATA                |CK|
            +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

             ID HID Report ID (BB)
             CK Checksum

    3.  Read (`03`)

             00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15
            +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
            |BB|02|03|00|00|00|00|00|00|00|00|00|00|00|00|40|
            +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
            |ID|     |             RESERVED              |CK|
            +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

             ID HID Report ID (BB)
             CK Checksum

             00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15
            +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
            |AA|02|03|20|20|20|20|20|20|20|20|20|20|20|20|D1|
            +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
            |ID|     |               DATA                |CK|
            +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

             ID HID Report ID (AA)
             CK Checksum

2.  Key Input (`03`)

         00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15
        +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
        |AA|03|00|00|00|00|00|00|00|00|00|00|00|00|00|53|
        +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
        |ID|  |XX|YY|            RESERVED            |CK|
        +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

         ID HID Report ID (AA)
         XX Key Code:
            00 Up
            01 Right
            02 Left
            03 Down
            04 Enter
            05 Cancel
         YY Key Event:
            00 Key Pressed
            01 Key Released
         CK Checksum

3.  Set Backlight (`07`)

         00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15
        +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
        |BB|07|00|00|00|00|00|00|00|00|00|00|00|00|00|3E|
        +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
        |ID|  |XX|             RESERVED              |CK|
        +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

         ID HID Report ID (BB)
         XX State:
            00 Off
            01 On
         CK Checksum

