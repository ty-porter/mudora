// Package sni wraps the SNI (Super Nintendo Interface) gRPC API with the
// small surface this tracker needs: device discovery and bulk memory reads.
//
// SNI: https://github.com/alttpo/sni
package sni

import (
	"context"
	"fmt"

	pb "github.com/ty-porter/mudora/internal/snipb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client is a thin wrapper around the SNI gRPC services.
type Client struct {
	conn    *grpc.ClientConn
	devices pb.DevicesClient
	memory  pb.DeviceMemoryClient
}

// Dial connects to an SNI server (default localhost:8191).
func Dial(ctx context.Context, addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dialing SNI: %w", err)
	}
	return &Client{
		conn:    conn,
		devices: pb.NewDevicesClient(conn),
		memory:  pb.NewDeviceMemoryClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// FirstDevice returns the URI of the first connected SNES device that
// supports memory reads, or an error if none are connected.
func (c *Client) FirstDevice(ctx context.Context) (string, error) {
	resp, err := c.devices.ListDevices(ctx, &pb.DevicesRequest{
		Kinds: nil, // all kinds
	})
	if err != nil {
		return "", fmt.Errorf("listing devices: %w", err)
	}
	for _, d := range resp.Devices {
		// TODO: filter on d.Capabilities containing ReadMemory.
		return d.Uri, nil
	}
	return "", fmt.Errorf("no SNES devices connected to SNI")
}

// ReadMemory reads size bytes from the device at the given FX Pak Pro
// address-space address. SNI translates address spaces per device, so
// FxPakPro space is the safe common currency (e.g. WRAM starts at 0xF50000).
func (c *Client) ReadMemory(ctx context.Context, deviceURI string, addr uint32, size uint32) ([]byte, error) {
	resp, err := c.memory.SingleRead(ctx, &pb.SingleReadMemoryRequest{
		Uri: deviceURI,
		Request: &pb.ReadMemoryRequest{
			RequestAddress:       addr,
			RequestAddressSpace:  pb.AddressSpace_FxPakPro,
			RequestMemoryMapping: pb.MemoryMapping_LoROM,
			Size:                 size,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("reading 0x%06x (%d bytes): %w", addr, size, err)
	}
	return resp.Response.Data, nil
}

// TODO: add MultiRead for batching several regions into one round trip
// (pb.MultiReadMemoryRequest) — important for keeping poll latency low on
// real hardware.
