package mid

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/testvergecloud/testApi/business/data/transaction"
	"github.com/testvergecloud/testApi/foundation/logger"
)

// ExecuteInTransaction starts a transaction around all the storage calls within
// the scope of the handler function.
// ExecuteInTransaction starts a transaction around all the storage calls within
// the scope of the handler function.
func ExecuteInTransaction(log *logger.Logger, bgn transaction.Beginner) gin.HandlerFunc {
	return func(c *gin.Context) {
		hasCommitted := false

		log.Info(c, "BEGIN TRANSACTION")
		tx, err := bgn.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("BEGIN TRANSACTION: %s", err)})
			c.Abort()
			return
		}

		defer func() {
			if !hasCommitted {
				log.Info(c, "ROLLBACK TRANSACTION")
			}

			if err := tx.Rollback(); err != nil {
				if errors.Is(err, sql.ErrTxDone) {
					return
				}
				log.Info(c, "ROLLBACK TRANSACTION", "ERROR", err)
			}
		}()

		c.Request = c.Request.WithContext(transaction.Set(c.Request.Context(), tx))

		c.Next()

		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("EXECUTE TRANSACTION: %s", c.Errors)})
			c.Abort()
			return
		}

		log.Info(c, "COMMIT TRANSACTION")
		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("COMMIT TRANSACTION: %s", err)})
			c.Abort()
			return
		}

		hasCommitted = true
	}
}
