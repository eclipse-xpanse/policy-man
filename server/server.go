/*
 * SPDX-License-Identifier: Apache-2.0
 * SPDX-FileCopyrightText: Huawei Inc.
 */

package server

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"github.com/eclipse-xpanse/policy-man/config"
	"github.com/eclipse-xpanse/policy-man/log"
	_ "github.com/eclipse-xpanse/policy-man/openapi/docs"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

func RunHTTPServer(ctx context.Context, cfg *config.Conf) error {
	server := &http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: router(cfg),
	}

	log.Info("HTTP server is running on " + cfg.Host + ":" + cfg.Port)
	if cfg.SSL.Enable {
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS10,
		}

		if tlsConfig.NextProtos == nil {
			tlsConfig.NextProtos = []string{"http/1.1"}
		}

		tlsConfig.Certificates = make([]tls.Certificate, 1)
		var err error
		if cfg.SSL.CertPath != "" && cfg.SSL.KeyPath != "" {
			tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(cfg.SSL.CertPath, cfg.SSL.KeyPath)
			if err != nil {
				log.Error("Failed to load https cert file: ", err)
				return err
			}
		} else if cfg.SSL.CertBase64 != "" && cfg.SSL.KeyBase64 != "" {
			cert, err := base64.StdEncoding.DecodeString(cfg.SSL.CertBase64)
			if err != nil {
				log.Error("base64 decode error:", err.Error())
				return err
			}
			key, err := base64.StdEncoding.DecodeString(cfg.SSL.KeyBase64)
			if err != nil {
				log.Error("base64 decode error:", err.Error())
				return err
			}
			if tlsConfig.Certificates[0], err = tls.X509KeyPair(cert, key); err != nil {
				log.Error("tls key pair error:", err.Error())
				return err
			}
		} else {
			return errors.New("missing https cert config")
		}

		server.TLSConfig = tlsConfig
	}

	return startServer(ctx, server, cfg)
}

func listenAndServe(ctx context.Context, s *http.Server, cfg *config.Conf) error {
	var g errgroup.Group
	g.Go(func() error {
		<-ctx.Done()
		timeout := time.Duration(cfg.ShutdownTimeout) * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return s.Shutdown(ctx)
	})
	g.Go(func() error {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	return g.Wait()
}

func listenAndServeTLS(ctx context.Context, s *http.Server, cfg *config.Conf) error {
	var g errgroup.Group
	g.Go(func() error {
		<-ctx.Done()
		timeout := time.Duration(cfg.ShutdownTimeout) * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return s.Shutdown(ctx)
	})
	g.Go(func() error {
		if err := s.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	return g.Wait()
}

func startServer(ctx context.Context, s *http.Server, cfg *config.Conf) error {
	if s.TLSConfig == nil {
		return listenAndServe(ctx, s, cfg)
	}

	return listenAndServeTLS(ctx, s, cfg)
}

func router(cfg *config.Conf) *gin.Engine {
	// set server mode
	gin.SetMode(cfg.Mode)

	r := gin.New()

	// Global middleware
	r.Use(logger.SetLogger(
		logger.WithUTC(true),
		logger.WithSkipPath([]string{}),
	))
	r.GET("/health", healthHandler)
	r.POST("/evaluate/policy", policyEvaluateHandler(cfg))
	r.POST("/evaluate/policies", policiesEvaluateHandler(cfg))
	r.POST("/validate/policies", PoliciesValidateHandler(cfg))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
