          config
             ↓
          runtime
        ↙     ↓     ↘
    logger   tracer   metrics
        ↘      ↓       ↙
             context
          ↙           ↘
    http middleware   grpc middleware
      ↓
    async
      ↓
    shutdown
